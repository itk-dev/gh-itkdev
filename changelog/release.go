package changelog

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

func getBranchName(release string) (string, error) {
	version := release
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}
	if !semver.IsValid(version) {
		return "", fmt.Errorf("invalid version: %s", release)
	}

	branchType := "release"
	// Determine if the release is a "release" or a "hotfix".
	//
	// A "hotfix" is a version that's not a prerelease and has a patch version different from 0 (cf. https://semver.org/).
	if semver.Prerelease(version) == "" && semver.MajorMinor(version)+".0" != version {
		branchType = "hotfix"
	}

	return branchType + "/" + release, nil
}

func createReleaseBranch(release string, base string) (string, error) {
	branch, err := getBranchName(release)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("git", "checkout", "-b", branch, base)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", output, err)
	}

	return branch, nil
}

func updateReleaseChangelog(changelog string, release string) (string, error) {
	headerPattern := regexp.MustCompile(`(?i)^\#+ +\[unreleased\]`)
	unreleasedHeaderIndex := -1
	// Match a URL on the form scheme://domain/user/repo
	linkUnreleasedPattern := regexp.MustCompile(`(?i)^\[unreleased\]: (?P<url>[a-z]+?://[^/]+/(?P<user>[^/]+)/(?P<repo>[^/]+))`)
	linkUnreleasedIndex := -1

	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(changelog))
	for scanner.Scan() {
		line := scanner.Text()
		if headerPattern.MatchString(line) {
			unreleasedHeaderIndex = len(lines)
		} else if linkUnreleasedPattern.MatchString(line) {
			linkUnreleasedIndex = len(lines)
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	if unreleasedHeaderIndex < 0 {
		return "", fmt.Errorf("cannot find \"Unreleased\"")
	} else if linkUnreleasedIndex < 0 {
		return "", fmt.Errorf("cannot find \"Unreleased\" link")
	} else if linkUnreleasedIndex < unreleasedHeaderIndex {
		return "", fmt.Errorf("\"Unreleased\" link must come after header")
	}

	// ----------------------------------------------------------------

	insertIndex := linkUnreleasedIndex + 1

	match := linkUnreleasedPattern.FindStringSubmatch(lines[linkUnreleasedIndex])
	url, err := url.Parse(match[1])
	if err != nil {
		return "", fmt.Errorf("invalid url: %s", match[1])
	}

	path := url.Path
	url.Path = fmt.Sprintf("%s/compare/%s...HEAD", path, release)
	unreleasedLink := fmt.Sprintf("[Unreleased]: %s", url)

	if insertIndex < len(lines) {
		line := lines[insertIndex]
		var previousRelease string
		if index := strings.LastIndex(line, "..."); index > -1 {
			previousRelease = line[index+3:]
		} else if index := strings.LastIndex(line, "/"); index > -1 {
			previousRelease = line[index+1:]
		} else {
			log.Fatalf("cannot find previous version from %q", line)
		}
		url.Path = fmt.Sprintf("%s/compare/%s...%s", path, previousRelease, release)
	} else {
		url.Path = fmt.Sprintf("%s/releases/tag/%s", path, release)
	}

	releaseLink := fmt.Sprintf("[%s]: %s", release, url)

	lines = slices.Concat(
		lines[0:insertIndex-1],
		[]string{
			unreleasedLink,
			releaseLink,
		},
		lines[insertIndex:],
	)

	// ----------------------------------------------------------------

	insertIndex = unreleasedHeaderIndex + 2

	item := fmt.Sprintf("## [%s] - %s\n", release, time.Now().Format("2006-01-02"))

	lines = slices.Concat(
		lines[0:insertIndex],
		[]string{item},
		lines[insertIndex:],
	)

	// ----------------------------------------------------------------

	return strings.Join(lines, "\n") + "\n", nil
}

func Release(release string, base string, name string, commit bool) {
	b, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	changelog := string(b)

	updatedChangelog, err := updateReleaseChangelog(changelog, release)
	if err != nil {
		log.Fatal(err)
	}

	branch, err := createReleaseBranch(release, base)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Branch %s created\n", branch)

	os.WriteFile(name, []byte(updatedChangelog), 0644)
	fmt.Printf("Updated changelog written to %s\n", name)

	if commit {
		gitCommit([]string{name}, fmt.Sprintf("Release %s", release))
		fmt.Println("Updated changelog committed")
	} else {
		gitDiff([]string{name})
	}
}
