package main

import (
	"encoding/json"
	"goBackupInspector/config"
	"goBackupInspector/diffs"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	var summaryFileName string
	if err := config.LoadConfiguration("configFile"); err != nil {
		log.Panic(err)
	}

	if err := config.CheckConfig(); err != nil {
		log.Panic(err)
	}

	summary, err := diffs.CreateFileDiff(config.Conf.BackupFilesDir, config.Conf.SysDir,
										 config.Conf.RelPathInBackup, config.Conf.DiffIgnoreFile)
	if err != nil {
		log.Panic(err)
	}
	if !summary.WithDifferences {
		return
	}

	summaryFileName = "diffSummary_" + time.Now().Format("2006_01_02") + ".json"
	file, _ := json.MarshalIndent(summary, "", "\t")
	_ = ioutil.WriteFile(summaryFileName, file, 0644)
	defer func() {
		if err := os.Remove(summaryFileName); err != nil {
			log.Fatal(err)
		}
	}()

	// todo error handling
	sendMail("There are unexpected diffs in your system. See attachment for details!", summaryFileName, *config.Conf)
	log.Print("summary sent to " + config.Conf.Mail.ReceiverAdr)
}
