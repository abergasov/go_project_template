package utils_test

import (
	"go_project_template/internal/utils"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeriveVersion(t *testing.T) {
	cases := map[string]string{
		// exact semver
		"v1.2.3": "v1.2.3",

		// prerelease
		"v1.2.3-rc1": "v1.2.3-rc1",

		// prerelease with +dirty
		"v0.0.1-rc4.0.20250821142859-6128ae7a7356+dirty": "v0.0.1-rc4",

		// pseudo base
		"v1.2.3-0.20240102112233-deadbeef": "v1.2.3",
		"v1.2.4-0.20240102112233-abcdef1":  "v1.2.4",

		// pure pseudo (v0.0.0 -> dev)
		"v0.0.0-20240102112233-deadbeef": "dev",
		"v0.0.0-20250821142859-deadbeef": "dev",

		// unknown8 preserved (special date)
		"v1.2.3-20250821-deadbeef": "v1.2.3-20250821-deadbeef",

		// devel fallback
		"(devel)": "dev",

		// build metadata stripped
		"v1.2.5+dirty":      "v1.2.5",
		"v1.2.5+build.meta": "v1.2.5",
	}

	for in, want := range cases {
		require.Equal(t, want, utils.DeriveVersion(in), "input=%q", in)
	}
}

func TestGetVersion(t *testing.T) {
	v := utils.GetVersion()

	require.NotEmpty(t, v, "version must not be empty")
	require.NotContains(t, v, "+", "version should be stripped of build metadata")

	// Accept:
	//   - "dev"
	//   - semver: vX.Y.Z (optionally with -prerelease tag)
	//   - non-standard tail will be  preserved: vX.Y.Z-YYYYMMDD-<hash> (special date format set)
	semverLike := regexp.MustCompile(`^dev$|^v\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?$|^v\d+\.\d+\.\d+-\d{8}-[0-9a-fA-F]+$`)
	require.True(t, semverLike.MatchString(v), "unexpected version format: %q", v)
}

func TestGetCommitHash(t *testing.T) {
	h := utils.GetCommitHash()
	// Either 7 chars or "unknown" which is 7 chars len too
	require.Len(t, h, 7, "CommitShort must be 7 chars or 'unknown'")
}
