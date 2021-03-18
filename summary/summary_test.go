package summary

import "testing"

func TestCompare(t *testing.T) {
	diffSummary := FileDiffSummary{LeftDir: "NotEqual"}
	other := FileDiffSummary{LeftDir: "../test/dump"}
	assertCompare(diffSummary, other, DiffLeftDir.i(), t)

	other.LeftDir = diffSummary.LeftDir
	diffSummary.RightDir = "RightDir"
	other.RightDir = "other"
	assertCompare(diffSummary, other, DiffRightDir.i(), t)

	other.RightDir = diffSummary.RightDir
	diffSummary.FilesNotInDir = map[string][]string{DirBackup: {"Not in dir"}}
	assertCompare(diffSummary, other, DiffFilesNotInDir.i(), t)

	other.FilesNotInDir = diffSummary.FilesNotInDir
	diffSummary.DirectoriesNotInDir = map[string][]string{"left": {"not in dir"}}
	assertCompare(diffSummary, other, DiffDirectoriesNotInDir.i(), t)

	other.DirectoriesNotInDir = diffSummary.DirectoriesNotInDir
	diffSummary.ComparedFiles = []string{"leftFile", "rightFile"}
	assertCompare(diffSummary, other, DiffComparedFiles.i(), t)

	other.ComparedFiles = diffSummary.ComparedFiles
	diffSummary.UnequalFiles = []FileTuple{{LeftFile: "left", RightFile: "right"}}
	assertCompare(diffSummary, other, DiffUnequalFiles.i(), t)

	other.UnequalFiles = diffSummary.UnequalFiles
	diffSummary.IgnoredElement = []IgnoredElement{{IgnoredElement: "file", CausedRule: "rule"}}
	assertCompare(diffSummary, other, DiffIgnoredFiles.i(), t)

	other.IgnoredElement = diffSummary.IgnoredElement
	diffSummary.WithDifferences = true
	assertCompare(diffSummary, other, DiffWithDifferences.i(), t)

	other.WithDifferences = diffSummary.WithDifferences
	diffSummary.BackupFileName = "backup file name"
	assertCompare(diffSummary, other, DiffBackupFiles.i(), t)

	other.BackupFileName = diffSummary.BackupFileName
	assertCompare(diffSummary, other, Equal.i(), t)
}

func assertCompare(sum, other FileDiffSummary, expected int, t *testing.T) {
	if sum.Compare(other) != expected {
		t.Fail()
	}
}
