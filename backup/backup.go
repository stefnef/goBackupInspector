package backup

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"regexp"
	"time"
)


func FindLastBackup(backupDir string) (backupName string, err error){
	regex, _ := regexp.Compile(".tar.gz")
	var lastModified time.Time

	files, err := ioutil.ReadDir(backupDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if regex.MatchString(file.Name()) {
			if file.ModTime().After(lastModified) {
				lastModified = file.ModTime()
				backupName = backupDir + file.Name()
			}
		}
	}

	if backupName == "" {
		err = errors.New("no backup file found")
	}

	return backupName, err
}

func createTMPDumpDir(dumpDir string) error {
	cmd := exec.Command("mkdir", dumpDir)
	err := cmd.Run()
	return err
}

func RemoveTMPDumpDir(dumpDir string) error {
	cmd := exec.Command("rm", "-r", dumpDir)
	err := cmd.Run()
	return err
}

func unzipBackup(backupFile, targetDir string) error {
	cmd := exec.Command("tar", "xf", backupFile, "-C", targetDir)
	return cmd.Run()
}

func Unpack(backupFile, targetDir string) error {
	if err := createTMPDumpDir(targetDir); err != nil {
		return err
	}

	if err := unzipBackup(backupFile, targetDir); err != nil {
		return err
	}
	return nil
}
