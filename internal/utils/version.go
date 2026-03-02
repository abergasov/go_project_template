package utils

import (
	"regexp"
	"runtime/debug"
	"strings"
)

const defaultDevVersion = "dev"
const defaultUnknownCommit = "unknown"

// GetVersion returns the application version derived from build info.
// Returns "dev" if version information is not available.
func GetVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok || bi == nil {
		return defaultDevVersion
	}
	v := DeriveVersion(bi.Main.Version)
	if v == "" {
		return defaultDevVersion
	}
	return v
}

// GetCommitHash returns the short commit hash (7 characters) from build info.
// Returns "unknown" if commit information is not available.
func GetCommitHash() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok || bi == nil {
		return defaultUnknownCommit
	}
	h := short(vcsRevision(bi))
	if h == "" {
		return defaultUnknownCommit
	}
	return h
}

// precompiled patterns (RE2-compatible: no non-capturing groups)
var (
	// vX.Y.Z or vX.Y.Z-<prerelease> -> kept as-is
	reExact = regexp.MustCompile(`^v\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?$`)

	// vX.Y.Z-<pre>.0.<14d>-<hash> -> vX.Y.Z-<pre>
	rePseudoWithPre = regexp.MustCompile(`^(v\d+\.\d+\.\d+-[0-9A-Za-z.-]+)\.0\.\d{14}-[0-9a-fA-F]+$`)

	// vX.Y.Z-0.<14d>-<hash> ->  vX.Y.Z
	rePseudoBase = regexp.MustCompile(`^(v\d+\.\d+\.\d+)-0\.\d{14}-[0-9a-fA-F]+$`)

	// vX.Y.Z-YYYYMMDD-<hash>  (kept as-is; not a Go pseudo)
	reUnknown8 = regexp.MustCompile(`^v\d+\.\d+\.\d+-\d{8}-[0-9a-fA-F]+$`)

	// v0.0.0-<14d>-<hash> -> dev
	rePurePseudo = regexp.MustCompile(`^v0\.0\.0-\d{14}-[0-9a-fA-F]+$`)
)

// DeriveVersion normalizes a version string from build info.
// Handles various version formats including semantic versions, pseudo-versions, and dev builds.
// Returns "dev" for unrecognized formats or pure pseudo-versions without a base tag.
func DeriveVersion(in string) string {
	v := stripBuildMeta(in)

	// explicit pseudo-forms first
	if m := rePseudoWithPre.FindStringSubmatch(v); m != nil {
		return m[1] // e.g., v1.2.3-rc4
	}
	if m := rePseudoBase.FindStringSubmatch(v); m != nil {
		return m[1] // e.g., v1.2.3
	}

	// Pure pseudo without a base tag
	if rePurePseudo.MatchString(v) {
		return defaultDevVersion
	}

	// Unknown 8-digit date tail is preserved
	if reUnknown8.MatchString(v) {
		return v
	}

	// Exact semver (incl. prerelease) that isn’t a pseudo
	if reExact.MatchString(v) {
		return v
	}

	// catch-all
	return defaultDevVersion
}

// -----------------------------
// Helpers
// -----------------------------
func stripBuildMeta(v string) string {
	if i := strings.IndexByte(v, '+'); i >= 0 {
		return v[:i]
	}
	return v
}

func vcsRevision(bi *debug.BuildInfo) string {
	for _, s := range bi.Settings {
		if s.Key == "vcs.revision" {
			return s.Value
		}
	}
	return ""
}

func short(h string) string {
	if len(h) >= 7 {
		return h[:7]
	}
	return h
}
