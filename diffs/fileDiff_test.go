package diffs

import (
	"goBackupInspector/summary"
	"regexp"
	"testing"
)

func TestPathWalk(t *testing.T) {
	filesExp := []string{"../test/sys/dir1/dir4/f5.txt", "../test/sys/dir1/dir4/f6.txt",
		"../test/sys/dir2/dir5/f7.txt", "../test/sys/dir2/dir6/f8.txt",
		"../test/sys/dir3/f1.txt", "../test/sys/dir3/f2.txt",
		"../test/sys/dir3/f3.txt", "../test/sys/dir4/f10.txt"}
	dirExp := []string{"../test/sys", "../test/sys/dir1", "../test/sys/dir1/dir4",
		"../test/sys/dir2", "../test/sys/dir2/dir5", "../test/sys/dir2/dir6",
		"../test/sys/dir3", "../test/sys/dir4"}
	if files, directories, err := filePathWalkDir("../test/sys"); err != nil {
		t.Error(err)
		t.FailNow()
	} else {
		if equal, actDiff, expDiff := assertEqual(files, filesExp); !equal {
			t.Errorf("files not as expected. Act '%s' vs Exp '%s'", actDiff, expDiff)
			t.Fail()
		}
		if equal, actDiff, expDiff := assertEqual(directories, dirExp); !equal {
			t.Errorf("directories not as expected. Act '%s' vs '%s'", actDiff, expDiff)
			t.Fail()
		}
	}
}

func assertEqual(act, exp []string) (isEquals bool, wrongAct, wrongExp string) {
	if len(act) != len(exp) {
		return false, "", ""
	}

	for idx, actElement := range act {
		if exp[idx] != actElement {
			return false, actElement, exp[idx]
		}
	}

	return true, "", ""
}

func TestIsIgnore(t *testing.T) {
	var ignoreList []*regexp.Regexp
	var err error
	var ignoreItems []string
	if ignoreItems, err = readDiffIgnore("../test/diffIgnore"); err != nil {
		t.Error(err)
		t.FailNow()
	}

	ignoreItemsExp := []string{"dir2/dir6/"}
	if isEqual, act, exp := assertEqual(ignoreItems, ignoreItemsExp); !isEqual {
		t.Errorf("ignoreItmes differ act: '%s' vs exp: '%s'", act, exp)
		t.FailNow()
	}

	if ignoreList, err = createIgnoreRegex(ignoreItems); err != nil {
		t.Error(err)
		t.FailNow()
	} else if len(ignoreList) != 1 {
		t.Errorf("len(ignoreList) = %d <> 1", len(ignoreList))
		t.FailNow()
	}

	if ignore, ruleName := isIgnored("dir2/dir6/", ignoreList); !ignore {
		t.Error("dir2/dir6/ not ignored")
	} else if ruleName != "dir2/dir6/" {
		t.Errorf("wrong rule name: '%s", ruleName)
	}
	if ignore, _ := isIgnored("../dir2/dir6/", ignoreList); !ignore {
		t.Error("../dir2/dir6 not ignored")
	}
	if ignore, _ := isIgnored("dir2/dir66/", ignoreList); ignore {
		t.Error("dir2/dir66/ is ignored")
	}
	if ignore, _ := isIgnored("dir2/dir6/555", ignoreList); !ignore {
		t.Error("dir2/dir6/555 is not ignored")
	}
}

func TestCreateFileDirChecklist(t *testing.T) {
	var files, directories []string
	var causes []summary.IgnoredElement
	filesExp := []string{"../test/sys/dir1/dir4/f5.txt", "../test/sys/dir1/dir4/f6.txt",
		"../test/sys/dir2/dir5/f7.txt",
		"../test/sys/dir3/f1.txt", "../test/sys/dir3/f2.txt",
		"../test/sys/dir3/f3.txt", "../test/sys/dir4/f10.txt"}
	directoriesExp := []string{"../test/sys", "../test/sys/dir1", "../test/sys/dir1/dir4",
		"../test/sys/dir2", "../test/sys/dir2/dir5",
		"../test/sys/dir3", "../test/sys/dir4"}

	if ignoreItems, err := readDiffIgnore("../test/diffIgnore"); err != nil {
		t.Error(err)
		t.FailNow()
	} else {
		if ignoreList, errCreate := createIgnoreRegex(ignoreItems); errCreate != nil {
			t.Error(errCreate)
			t.FailNow()
		} else {
			files, directories, causes = createFileDirChecklist("../test/sys", ignoreList)
		}
	}
	if isEqual, act, exp := assertEqual(files, filesExp); !isEqual {
		t.Errorf("files not as expected. Act: '%s' vs exp: '%s'", act, exp)
		t.Fail()
	}
	if isEqual, act, exp := assertEqual(directories, directoriesExp); !isEqual {
		t.Errorf("directories not as expected. Act: '%s' vs exp: '%s'", act, exp)
		t.FailNow()
	}
	if len(causes) != 2 {
		t.Errorf("wrong number of caused ignore rules found: %d <> 2", len(causes))
	}
}

func TestFindCuts(t *testing.T) {
	left := []string{"left/file1", "left/dir1/file2"}

	right := []string{"right/file0", "right/file1"}
	var filesToCompare []summary.FileTuple
	compareExp := summary.FileTuple{LeftFile: "left/file1", RightFile: "right/file1"}
	var filesNotInRightDir []string

	filesToCompare, filesNotInRightDir = findCuts("left", left, "right", right)
	if len(filesToCompare) != 1 {
		t.Error("wrong nr of compare")
	} else if filesToCompare[0] != compareExp {
		t.Error("wrong file to compare")
	}
	if len(filesNotInRightDir) != 1 {
		t.Error("wrong nr of compare")
	} else if filesNotInRightDir[0] != "left/dir1/file2" {
		t.Errorf("wrong file: act '%s' vs exp 'left/dir1/file2'", filesNotInRightDir[0])
	}
}

func TestFileDiff(t *testing.T) {
	var diffSummary summary.FileDiffSummary
	summaryExp := summary.FileDiffSummary{
		LeftDir:  "../test/dump",
		RightDir: "../test/sys",
		FilesNotInDir: map[string][]string{
			"../test/dump": {"../test/sys/dir2/dir5/f7.txt", "../test/sys/dir4/f10.txt"},
			"../test/sys":  {"../test/dump/dir3/f4.txt"}},
		DirectoriesNotInDir: map[string][]string{
			"../test/dump": {"../test/sys/dir4"},
			"../test/sys":  []string{},
		},
		ComparedFiles: []summary.FileTuple{
			{LeftFile: "../test/dump/dir1/dir4/f5.txt", RightFile: "../test/sys/dir1/dir4/f5.txt"},
			{LeftFile: "../test/dump/dir1/dir4/f6.txt", RightFile: "../test/sys/dir1/dir4/f6.txt"},
			{LeftFile: "../test/dump/dir3/f1.txt", RightFile: "../test/sys/dir3/f1.txt"},
			{LeftFile: "../test/dump/dir3/f2.txt", RightFile: "../test/sys/dir3/f2.txt"},
			{LeftFile: "../test/dump/dir3/f3.txt", RightFile: "../test/sys/dir3/f3.txt"},
		},
		UnequalFiles: []summary.FileTuple{
			{LeftFile: "../test/dump/dir1/dir4/f6.txt", RightFile: "../test/sys/dir1/dir4/f6.txt"},
		},
		IgnoredElement: []summary.IgnoredElement{
			{IgnoredElement: "../test/dump/dir2/dir6/f8.txt", CausedRule: "dir2/dir6/"},
			{IgnoredElement: "../test/dump/dir2/dir6", CausedRule: "dir2/dir6/"},
			{IgnoredElement: "../test/sys/dir2/dir6/f8.txt", CausedRule: "dir2/dir6/"},
			{IgnoredElement: "../test/sys/dir2/dir6", CausedRule: "dir2/dir6/"},
		},
		WithDifferences: true,
	}
	diffSummary = fileDiff("../test/dump", "../test/sys", "../test/diffIgnore")
	assertSummary(diffSummary, summaryExp, t)
}

func TestCreateFileDiff(t *testing.T) {
	_, err := CreateFileDiff("../test/backupFiles/", "../test/sys", "dump", "../test/diffIgnore")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func assertSummary(act, exp summary.FileDiffSummary, t *testing.T) {
	if cmp := act.Compare(exp); cmp != 0 {
		t.Errorf("summary not as expected, code %d", cmp)
		t.Fail()
	}
}
