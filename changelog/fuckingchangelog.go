package changelog

import (
	"bufio"
	"fmt"
	"iter"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
	"text/template"
)

// https://stackoverflow.com/a/71624929
func Map[T, U any](seq iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for a := range seq {
			if !yield(f(a)) {
				return
			}
		}
	}
}

func addPullRequest(changelog string, pr pullRequest, itemTemplate string) (string, error) {
	unreleasedHeaders := []string{"Unreleased", "Under udvikling"}
	var quotedHeaders []string
	for header := range Map(slices.Values(unreleasedHeaders), regexp.QuoteMeta) {
		quotedHeaders = append(quotedHeaders, header)
	}
	headerPattern := regexp.MustCompile(fmt.Sprintf(`(?i)^\#+ +\[(%s)\]`, strings.Join(quotedHeaders, "|")))

	unreleasedHeaderIndex := -1
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(changelog))
	for scanner.Scan() {
		line := scanner.Text()
		if headerPattern.MatchString(line) {
			unreleasedHeaderIndex = len(lines)
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	if unreleasedHeaderIndex < 0 {
		return "", fmt.Errorf("cannot find %s header", strings.Join(unreleasedHeaders, "/"))
	}

	// Make sure that we have a blank line after the header
	if !(unreleasedHeaderIndex < len(lines)-2) {
		lines = append(lines, "")
	}
	insertIndex := unreleasedHeaderIndex + 2
	tmpl, err := template.New("item").Parse(itemTemplate)
	if err != nil {
		return "", fmt.Errorf("cannot parse item template: %s", err)
	}
	var builder strings.Builder
	tmpl.Execute(&builder, pr)
	item := builder.String()

	// If content right after insertion point is a header or a link, we insert a blank line
	if insertIndex < len(lines) && regexp.MustCompile(`^[\[#]`).MatchString(lines[insertIndex]) {
		item += "\n"
	}

	lines = slices.Concat(
		lines[0:insertIndex],
		[]string{item},
		lines[insertIndex:],
	)

	return strings.Join(lines, "\n") + "\n", nil
}

func FuckingChangelog(name string, itemTemplate string) {
	b, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	changelog := string(b)

	if isChanged(name) {
		log.Fatalf("File %s is changed", name)
	}

	pr, err := getPullRequest()
	if err != nil {
		log.Fatalf("error getting pull request: %s\n", err)
	}

	updatedChangelog, err := addPullRequest(changelog, pr, itemTemplate)
	if err != nil {
		log.Fatalf("error adding pull request: %s\n", err)
	}

	os.WriteFile(name, []byte(updatedChangelog), 0644)
	fmt.Printf("Updated changelog written to %s\n", name)

	gitDiff([]string{name})
}
