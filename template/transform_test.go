package template

import (
	"goBackupInspector/summary"
	"os"
	"testing"
	"time"
)

func TestTransform(t *testing.T) {
	var fileName string
	var err error

	sum := createExampleSummary()

	if fileName, err = CreateHTMLFile(sum, ""); err != nil {
		t.Error(err)
	}
	if f, err := os.Open(fileName); err != nil {
		t.Error(err)
		t.Fail()
	} else {
		if errRemove := os.Remove(f.Name()); errRemove != nil {
			t.Error(errRemove)
		}
	}
}

func createExampleSummary() *summary.FileDiffSummary {
	return &summary.FileDiffSummary{
		Date:      time.Now(),
		BackupDir: "BackupDir",
		SystemDir: "SystemDir",
		FilesNotInDir: map[string][]string{
			summary.DirBackup: {"../File1", "../File2"},
			summary.DirSystem: {"../../Right3"}},
		DirectoriesNotInDir: map[string][]string{
			summary.DirBackup: {"../Dir1", "../Dir2"},
			summary.DirSystem: {"../../Dir3"}},
		ComparedFiles: []string{"Left1", "Left2"},
		IgnoredElements: []summary.IgnoredElement{
			{IgnoredElement: "Ign1", CausedRule: "Cause1"}, {IgnoredElement: "Ignore2", CausedRule: "C2"},
		},
		UnequalFiles: []summary.FileTuple{
			{BackupFile: "Left1", SystemFile: "Right1"}, {BackupFile: "Left2", SystemFile: "Right2"},
			{BackupFile: "Left3", SystemFile: "Right3"},
		},
		WithDifferences: true,
		BackupFileName:  "BackupFile.tar.gz",
	}
}
