package changelog

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddRelease(t *testing.T) {
	testCases := []struct {
		changelog string
		release   string
		expected  string
		errorText string
	}{
		{
			`## [Unreleased]

[Unreleased]: https://git.example.com
`,
			"v0.0.0",
			`## [Unreleased]

## [v0.0.0] - %TODAY%

[Unreleased]: https://git.example.com/compare/v0.0.0...HEAD
[v0.0.0]: https://git.example.com/releases/tag/v0.0.0
`,
			"cannot find \"Unreleased\" link",
		},

		{
			`## [Unreleased]

[Unreleased]: https://git.example.com/user/repo/
`,
			"v0.0.0",
			`## [Unreleased]

## [v0.0.0] - %TODAY%

[Unreleased]: https://git.example.com/user/repo/compare/v0.0.0...HEAD
[v0.0.0]: https://git.example.com/user/repo/releases/tag/v0.0.0
`,
			"",
		},

		{
			`## [Unreleased]

* [PR-42](https://git.example.com/user/repo/pr/42)
  Added the meaning

## [v0.0.0] - 2001-01-01

[Unreleased]: https://git.example.com/user/repo/compare/v0.0.0...HEAD
[v0.0.0]: https://git.example.com/user/repo/releases/tag/v0.0.0
`,
			"v0.1.0",
			`## [Unreleased]

## [v0.1.0] - %TODAY%

* [PR-42](https://git.example.com/user/repo/pr/42)
  Added the meaning

## [v0.0.0] - 2001-01-01

[Unreleased]: https://git.example.com/user/repo/compare/v0.1.0...HEAD
[v0.1.0]: https://git.example.com/user/repo/compare/v0.0.0...v0.1.0
[v0.0.0]: https://git.example.com/user/repo/releases/tag/v0.0.0
`,
			"",
		},

		{
			`## [Unreleased]

* [PR-42](https://git.example.com/user/repo/pr/42)
  Added the meaning

## [v0.1.0] - 2002-01-01

## [v0.0.0] - 2001-01-01

[Unreleased]: https://git.example.com/user/repo/compare/v0.1.0...HEAD
[v0.1.0]: https://git.example.com/user/repo/compare/v0.0.0...v0.1.0
[v0.0.0]: https://git.example.com/user/repo/releases/tag/v0.0.0
`,
			"v0.1.1",
			`## [Unreleased]

## [v0.1.1] - %TODAY%

* [PR-42](https://git.example.com/user/repo/pr/42)
  Added the meaning

## [v0.1.0] - 2002-01-01

## [v0.0.0] - 2001-01-01

[Unreleased]: https://git.example.com/user/repo/compare/v0.1.1...HEAD
[v0.1.1]: https://git.example.com/user/repo/compare/v0.1.0...v0.1.1
[v0.1.0]: https://git.example.com/user/repo/compare/v0.0.0...v0.1.0
[v0.0.0]: https://git.example.com/user/repo/releases/tag/v0.0.0
`,
			"",
		},
	}

	for _, testCase := range testCases {
		actual, err := updateReleaseChangelog(testCase.changelog, testCase.release)

		if testCase.errorText != "" {
			assert.NotNil(t, err)
			assert.EqualError(t, err, testCase.errorText)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, strings.ReplaceAll(testCase.expected, "%TODAY%", time.Now().Format("2006-01-02")), actual)
		}
	}
}
