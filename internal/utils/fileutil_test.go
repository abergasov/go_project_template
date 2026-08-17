package utils_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc64"
	"os"
	"path/filepath"
	"testing"

	"go_project_template/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestCleanupFolders_CreateAndClean(t *testing.T) {
	tmp := t.TempDir()
	// target dirs
	a := filepath.Join(tmp, "a")
	b := filepath.Join(tmp, "b")

	// pre-populate with files and subdirs (nested including)
	require.NoError(t, os.MkdirAll(filepath.Join(a, "sub"), 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(a, "keep.txt"), []byte("keep"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(a, "kill.txt"), []byte("kill"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(a, "sub", "nested.txt"), []byte("nested"), 0o600))
	// b does not exist yet; CleanupFolders should create it

	// Ignore keep.txt only; subdir and kill.txt should be removed
	require.NoError(t, utils.CleanupFolders([]string{a, b}, "keep.txt"))

	// a exists, b created
	_, err := os.Stat(a)
	require.NoError(t, err)
	_, err = os.Stat(b)
	require.NoError(t, err)

	// keep.txt remains
	_, err = os.Stat(filepath.Join(a, "keep.txt"))
	require.NoError(t, err)

	// kill.txt and subdir removed
	_, err = os.Stat(filepath.Join(a, "kill.txt"))
	require.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(a, "sub"))
	require.True(t, os.IsNotExist(err))
}

func TestCleanupFolders_IgnoreDoesNotRemoveDirectoriesWithSameNameAsIgnoredFile(t *testing.T) {
	tmp := t.TempDir()
	dir := filepath.Join(tmp, "dir")
	require.NoError(t, os.MkdirAll(dir, 0o750))
	// create a directory that has the same name as an ignored "file" entry
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "keepme"), 0o750))
	// create another file to be removed
	require.NoError(t, os.WriteFile(filepath.Join(dir, "deleteme"), []byte("x"), 0o600))

	// only exact entry name is ignored; the directory named "keepme" shall not be removed
	require.NoError(t, utils.CleanupFolders([]string{dir}, "keepme"))

	_, err := os.Stat(filepath.Join(dir, "keepme"))
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(dir, "deleteme"))
	require.True(t, os.IsNotExist(err))
}

func TestAtomicallySaveToFile_PrefixChecksumFormat(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "blob.bin")
	payload := []byte("payload-data-123")

	require.NoError(t, utils.AtomicallySaveToFile(target, payload))

	// read raw file and verify first utils.ChecksumSize bytes == crc64(remaining)
	// #nosec G304 -- target is generated in a safe test TempDir; not user-controlled
	raw, err := os.ReadFile(target)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(raw), utils.ChecksumSize)

	wantCRC := crc64.Checksum(payload, crc64.MakeTable(crc64.ISO))
	gotCRC := binary.LittleEndian.Uint64(raw[:utils.ChecksumSize])
	require.Equal(t, wantCRC, gotCRC)
	require.True(t, bytes.Equal(payload, raw[utils.ChecksumSize:]))
}

func TestAtomicallySaveToFile_EmptyPayload(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "empty.bin")
	require.NoError(t, utils.AtomicallySaveToFile(target, nil))

	// should still be utils.ChecksumSize bytes (crc of empty slice)
	// #nosec G304 -- target is generated in a safe test TempDir; not user-controlled
	raw, err := os.ReadFile(target)
	require.NoError(t, err)
	require.Equal(t, utils.ChecksumSize, len(raw))

	data, err := utils.LoadFromFile(target)
	require.NoError(t, err)
	require.Empty(t, data)
}

func TestAtomicallySaveToFile_PreserveExistingMode(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "mode.bin")

	// create file with a specific mode, then save over it
	require.NoError(t, os.WriteFile(target, []byte("tmp"), 0o600))
	require.NoError(t, os.Chmod(target, 0o600))

	require.NoError(t, utils.AtomicallySaveToFile(target, []byte("content")))
	st, err := os.Stat(target)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o600), st.Mode().Perm())
}

func TestLoadFromFile_PathIsCleaned(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "data.bin")
	payload := []byte("abc")
	require.NoError(t, utils.AtomicallySaveToFile(target, payload))

	// use a path with ../ segments; function calls filepath.Clean internally
	cleaned := filepath.Join(tmp, "x", "..", "data.bin")
	data, err := utils.LoadFromFile(cleaned)
	require.NoError(t, err)
	require.Equal(t, payload, data)
}

// ---------------
// --- helpers ---
// ---------------

// compile-time interface assertion
var _ fmt.Stringer = stringer("")

type stringer string

func (s stringer) String() string { return string(s) }
