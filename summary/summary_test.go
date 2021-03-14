package summary

import "testing"

func TestCompare(t *testing.T) {
	diffSummary := FileDiffSummary{LeftDir: "NotEqual"}
	other := FileDiffSummary{LeftDir: "../test/dump"}
	assertCompare(diffSummary, other, LeftDir.i(), t)

	other.LeftDir = diffSummary.LeftDir
	diffSummary.RightDir = "RightDir"
	other.RightDir = "other"
	assertCompare(diffSummary, other, RightDir.i(), t)

	other.RightDir = diffSummary.RightDir
	diffSummary.FilesNotInDir = map[string][]string{"left": {"Not in dir"}}
	assertCompare(diffSummary, other, FilesNotInDir.i(), t)

	other.FilesNotInDir = diffSummary.FilesNotInDir
	diffSummary.DirectoriesNotInDir = map[string][]string{"left": {"not in dir"}}
	assertCompare(diffSummary, other, DirectoriesNotInDir.i(), t)

	other.DirectoriesNotInDir = diffSummary.DirectoriesNotInDir
	diffSummary.ComparedFiles = []FileTuple{{LeftFile: "left", RightFile: "right"}}
	assertCompare(diffSummary, other, ComparedFiles.i(), t)

	other.ComparedFiles = diffSummary.ComparedFiles
	diffSummary.UnequalFiles = []FileTuple{{LeftFile: "left", RightFile: "right"}}
	assertCompare(diffSummary, other, UnequalFiles.i(), t)

	other.UnequalFiles = diffSummary.UnequalFiles
	diffSummary.IgnoredElement = []IgnoredElement{{IgnoredElement: "file", CausedRule: "rule"}}
	assertCompare(diffSummary, other, IgnoredFiles.i(), t)

	other.IgnoredElement = diffSummary.IgnoredElement
	diffSummary.WithDifferences = true
	assertCompare(diffSummary, other, WithDifferences.i(), t)

	other.WithDifferences = diffSummary.WithDifferences
	diffSummary.BackupFileName = "backup file name"
	assertCompare(diffSummary, other, BackupFiles.i(), t)

	other.BackupFileName = diffSummary.BackupFileName
	assertCompare(diffSummary, other, Equal.i(), t)
}

func assertCompare(sum, other FileDiffSummary, expected int, t *testing.T) {
	if sum.Compare(other) != expected {
		t.Fail()
	}
}
