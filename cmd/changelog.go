package cmd

import (
	"fmt"
	"os/exec"

	"github.com/itk-dev/gh-itkdev/changelog"
	"github.com/spf13/cobra"
)

// changelogCmd represents the changelog command
var (
	create bool

	fuckingChangelog        bool
	pullRequestItemTemplate string = func() string {
		format, _ := changelog.DetectPullRequestEntryFormat(changelogName)

		return format
	}()

	release    string
	baseBranch string = func() string {
		branches := []string{"develop", "main", "master"}
		for _, branch := range branches {
			cmd := exec.Command("git", "rev-parse", "--verify", "--branch", branch)
			if _, err := cmd.CombinedOutput(); err == nil {
				return branch
			}
		}
		// Fallback if no other suitable branch is found.
		return branches[0]
	}()
	commit bool = false

	changelogName string = "CHANGELOG.md"

	changelogCmd = &cobra.Command{
		Use:   "changelog",
		Short: "Update changelog",
		Run: func(cmd *cobra.Command, args []string) {
			if create {
				changelog.Create(changelogName)
			} else if fuckingChangelog {
				changelog.FuckingChangelog(changelogName, pullRequestItemTemplate)
			} else if release != "" {
				changelog.Release(release, baseBranch, changelogName, commit)
			} else {
				cmd.Usage()
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(changelogCmd)

	changelogCmd.Flags().BoolVarP(&create, "create", "", false, fmt.Sprintf("create a changelog (%q) if it does not exist", changelogName))

	changelogCmd.Flags().BoolVarP(&fuckingChangelog, "fucking-changelog", "", false, "add missing pull request entry to changelog")
	changelogCmd.Flags().StringVarP(&pullRequestItemTemplate, "item-template", "", pullRequestItemTemplate, "pull request item template")

	changelogCmd.Flags().StringVarP(&release, "release", "", "", "create a release branch with updated changelog")
	changelogCmd.Flags().StringVarP(&baseBranch, "base", "", baseBranch, "base branch for release")
	changelogCmd.Flags().BoolVarP(&commit, "commit", "", commit, "commit changes")

	changelogCmd.Flags().StringVarP(&changelogName, "changelog", "", changelogName, "changelog name")
}
