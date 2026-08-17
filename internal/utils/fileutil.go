package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc64"
	stdio "io"
	"os"
	"path/filepath"
	"slices"
)

// ChecksumSize is the length (in bytes) of the CRC64 checksum header
const ChecksumSize = 8

// CleanupFolders ensures that each directory in dirs exists and removes all entries except in ignoreFiles
func CleanupFolders(dirs []string, ignoreFiles ...string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return fmt.Errorf("cannot create directory %q: %w", dir, err)
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("cannot read directory %q: %w", dir, err)
		}
		for _, entry := range entries {
			if slices.Contains(ignoreFiles, entry.Name()) {
				continue
			}
			entryPath := filepath.Join(dir, entry.Name())
			if err = os.RemoveAll(entryPath); err != nil {
				return fmt.Errorf("cannot remove %q: %w", entryPath, err)
			}
		}
	}
	return nil
}

// AtomicallySaveToFile saves given data to the given file atomically
// Appends a CRC64 checksum to the data before writing
// File is either fully updated or not updated at all -> done by writing to a temporary file and renaming after
func AtomicallySaveToFile(fileName string, data []byte) (err error) {
	checkSum := crc64.New(crc64.MakeTable(crc64.ISO))
	if _, err = checkSum.Write(data); err != nil {
		return fmt.Errorf("cannot calculate checksum: %w", err)
	}
	sum := make([]byte, ChecksumSize)
	binary.LittleEndian.PutUint64(sum, checkSum.Sum64())
	result := make([]byte, ChecksumSize+len(data))
	copy(result[:ChecksumSize], sum)
	copy(result[ChecksumSize:], data)

	dir, base := filepath.Split(fileName)
	if dir == "" {
		dir = "."
	}
	f, err := os.CreateTemp(dir, base)
	if err != nil {
		return fmt.Errorf("cannot create temp file: %w", err)
	}
	// always try to close, even if we return early
	defer f.Close() //nolint:errcheck

	tempPath := f.Name()

	// write + flush
	if _, err = stdio.Copy(f, bytes.NewReader(result)); err != nil {
		return fmt.Errorf("cannot write data to tempfile %q: %w", tempPath, err)
	}
	if err = f.Sync(); err != nil {
		return fmt.Errorf("can't flush tempfile %q: %w", tempPath, err)
	}
	if err = f.Close(); err != nil { // explicit close before rename (Windows-friendly)
		return fmt.Errorf("can't close tempfile %q: %w", tempPath, err)
	}

	// preserve mode if destination exists
	if destInfo, statErr := os.Stat(fileName); statErr == nil {
		if srcInfo, err := os.Stat(tempPath); err == nil && srcInfo.Mode() != destInfo.Mode() {
			if err = os.Chmod(tempPath, destInfo.Mode()); err != nil {
				return fmt.Errorf("can't set filemode on tempfile %q: %w", tempPath, err)
			}
		}
	} else if !os.IsNotExist(statErr) {
		return statErr
	}

	// atomic replace
	if err = os.Rename(tempPath, fileName); err != nil {
		_ = os.Remove(tempPath) // clean up temp file immediately; no flag logic
		return fmt.Errorf("cannot replace %q with tempfile %q: %w", fileName, tempPath, err)
	}

	// fsync directory to make rename durable
	// #nosec G304 -- safe: dir is derived from controlled file path
	dirFD, err := os.Open(dir)
	if err == nil {
		if syncErr := dirFD.Sync(); syncErr != nil {
			// not fatal everywhere
			_ = dirFD.Close()
			return fmt.Errorf("directory fsync failed for %q: %w", dir, syncErr)
		}
		_ = dirFD.Close()
	} // if opening dir fails fsync is skipped, may be logged

	return nil
}

// LoadFromFile loads data from the given file and verifies the checksum
// Returns data w/o checksum
func LoadFromFile(path string) (data []byte, err error) {
	path = filepath.Clean(path)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %q does not exist: %w", path, os.ErrNotExist)
	}
	// #nosec G304 -- path is caller-controlled within the application; sanitized above
	r, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer r.Close() //nolint:errcheck

	data, err = stdio.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if len(data) < ChecksumSize {
		return nil, errors.New("file is too short")
	}
	fileCrc := binary.LittleEndian.Uint64(data[:ChecksumSize])
	dataCrc := crc64.Checksum(data[ChecksumSize:], crc64.MakeTable(crc64.ISO))
	if fileCrc != dataCrc {
		return nil, fmt.Errorf("checksum mismatch: %x != %x", fileCrc, dataCrc)
	}
	return data[ChecksumSize:], nil
}
