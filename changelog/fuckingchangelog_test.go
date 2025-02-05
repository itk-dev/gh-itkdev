package changelog

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuckingChangelog(t *testing.T) {
	defaultItemTemplate := `* [PR-{{ .Number }}]({{ .Url }})
  {{ .Title }}`
	testCases := []struct {
		changelog    string
		pullRequest  pullRequest
		itemTemplate string
		expected     string
		expectedErr  error
	}{
		{
			`## [Unreleased]
`,
			pullRequest{
				Number: 87,
				Title:  "Test",
				Url:    "https://example.com/pr/87",
			},
			defaultItemTemplate,
			`## [Unreleased]

* [PR-87](https://example.com/pr/87)
  Test
`,
			nil,
		},

		{
			`## [Unreleased]

[Unreleased]: https://example.com/
`,
			pullRequest{
				Number: 87,
				Title:  "Test",
				Url:    "https://example.com/pr/87",
			},
			defaultItemTemplate,
			`## [Unreleased]

* [PR-87](https://example.com/pr/87)
  Test

[Unreleased]: https://example.com/
`,
			nil,
		},

		{
			`## [Unreleased]

* [PR-42](https://example.com/pr/42)
  Added the meaning
`,
			pullRequest{
				Number: 87,
				Title:  "Test",
				Url:    "https://example.com/pr/87",
			},
			defaultItemTemplate,
			`## [Unreleased]

* [PR-87](https://example.com/pr/87)
  Test
* [PR-42](https://example.com/pr/42)
  Added the meaning
`,
			nil,
		},

		{
			`## [Unreleased]

- [#42](https://example.com/pr/42): Added the meaning
`,
			pullRequest{
				Number: 87,
				Title:  "Test",
				Url:    "https://example.com/pr/87",
			},
			`- [#{{ .Number }}]({{ .Url }}): {{ .Title }}`,
			`## [Unreleased]

- [#87](https://example.com/pr/87): Test
- [#42](https://example.com/pr/42): Added the meaning
`,
			nil,
		},

		{
			`## [Unreleased]

## [v0.0.0] - 2001-01-01

- [#42](https://example.com/pr/42): Added the meaning
`,
			pullRequest{
				Number: 87,
				Title:  "Test",
				Url:    "https://example.com/pr/87",
			},
			`- [#{{ .Number }}]({{ .Url }}): {{ .Title }}`,
			`## [Unreleased]

- [#87](https://example.com/pr/87): Test

## [v0.0.0] - 2001-01-01

- [#42](https://example.com/pr/42): Added the meaning
`,
			nil,
		},

		{
			`## [Under udvikling]
`,
			pullRequest{
				Number: 87,
				Title:  "Test",
				Url:    "https://example.com/pr/87",
			},
			defaultItemTemplate,
			`## [Under udvikling]

* [PR-87](https://example.com/pr/87)
  Test
`,
			nil,
		},

		{
			`## [under udviklinG]
`,
			pullRequest{
				Number: 87,
				Title:  "Test",
				Url:    "https://example.com/pr/87",
			},
			defaultItemTemplate,
			`## [under udviklinG]

* [PR-87](https://example.com/pr/87)
  Test
`,
			nil,
		},

		{
			`## [We're still working on it]
`,
			pullRequest{
				Number: 87,
				Title:  "Test",
				Url:    "https://example.com/pr/87",
			},
			defaultItemTemplate,
			"",
			fmt.Errorf("cannot find Unreleased/Under udvikling header"),
		},
	}

	for _, testCase := range testCases {
		actual, err := addPullRequest(testCase.changelog, testCase.pullRequest, testCase.itemTemplate)
		if err != nil {
			if assert.NotNil(t, testCase.expectedErr) {
				assert.Equal(t, testCase.expectedErr.Error(), err.Error())
			}
		} else {
			assert.Equal(t, testCase.expected, actual)
		}
	}
}
