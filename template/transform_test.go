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
		Date:     time.Now(),
		LeftDir:  "LeftDir",
		RightDir: "RightDir",
		FilesNotInDir: map[string][]string{
			"LeftDir":  {"../File1", "../File2"},
			"RightDir": {"../../Right3"}},
		DirectoriesNotInDir: map[string][]string{
			"LeftDir":  {"../Dir1", "../Dir2"},
			"RightDir": {"../../Dir3"}},
		ComparedFiles: []summary.FileTuple{
			{LeftFile: "Left1", RightFile: "Right1"}, {LeftFile: "Left2", RightFile: "Right2"},
		},
		IgnoredElement: []summary.IgnoredElement{
			{IgnoredElement: "Ign1", CausedRule: "Cause1"}, {IgnoredElement: "Ignore2", CausedRule: "C2"},
		},
		UnequalFiles: []summary.FileTuple{
			{LeftFile: "Left1", RightFile: "Right1"}, {LeftFile: "Left2", RightFile: "Right2"},
			{LeftFile: "Left3", RightFile: "Right3"},
		},
		WithDifferences: true,
		BackupFileName:  "BackupFile.tar.gz",
	}
}
