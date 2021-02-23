package backup

import (
	"goBackupInspector/config"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	config.LoadConfiguration("../test/configTestFile")
	if config.Conf == nil {
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestUnpackBackupDir(t *testing.T) {
	if _, err := os.Stat("../test/backupFiles/backupDUMP.tar.gz"); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err := Unpack("../test/backupFiles/backupDUMP.tar.gz", "../test/tmpDUMP"); err != nil {
		t.Error(err)
		t.Fail()
	}
	_, err := os.Stat("../test/tmpDUMP")
	if os.IsNotExist(err) {
		t.Error(err)
	}

	defer deferRemoveTMPDumpDir(t)

}

func TestFindLastBackup(t *testing.T) {
	files := []string{"../test/backup_1.tar.gz", "../test/backup_2.tar.gz", "../test/backup_3.tar.gz"}
	for _,file := range files {
		if err := createZipFile(file, "../test/sys", t); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
	defer deferRemoveZipFiles(files, t)

	if lastFile, err := FindLastBackup("../test/"); err != nil {
		t.Error(err)
	} else if lastFile != files[2] {
		t.Errorf("wrong file found. Found: '%s', expected: '%s'", lastFile, files[2])
	}
}

func TestCreateTMPDumpDir(t *testing.T) {
	if err := createTMPDumpDir("../test/tmpDUMP"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err := os.Stat("../test/tmpDUMP")
	if os.IsNotExist(err) {
		t.Error(err)
	}

	if err := RemoveTMPDumpDir("../test/tmpDUMP"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = os.Stat("../test/tmpDUMP")
	if !os.IsNotExist(err) {
		t.Error(err)
	}
}

func TestUnzipBackup(t *testing.T) {
	var err error
	if err = createTMPDumpDir("../test/tmpDUMP"); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err = unzipBackup("../test/backupFiles/backupDUMP.tar.gz", "../test/tmpDUMP"); err != nil {
		t.Error(err)
		t.Fail()
	}

	_, err = os.Stat("../test/tmpDUMP")
	if os.IsNotExist(err) {
		t.Error(err)
	}

	defer deferRemoveTMPDumpDir(t)
}

func deferRemoveTMPDumpDir(t *testing.T) {
	err := RemoveTMPDumpDir("../test/tmpDUMP")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func deferRemoveZipFiles(files []string, t *testing.T) {
	for _, file := range files {
		err := RemoveTMPDumpDir(file)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}
}

func createZipFile(fileName, backupDir string, t *testing.T) error {
	cmd := exec.Command("tar", "cf", fileName, backupDir)
	return cmd.Run()
}