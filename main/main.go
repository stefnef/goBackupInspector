package main

import (
	"goBackupInspector/config"
	"goBackupInspector/diffs"
	"goBackupInspector/template"
	"log"
	"os"
)

func main() {
	var summaryFileName string
	if len(os.Args) != 2 {
		if err := config.LoadConfiguration("configFile"); err != nil {
			log.Panic(err)
		}
	} else {
		if err := config.LoadConfiguration(os.Args[1]); err != nil {
			log.Panic(err)
		}
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

	if config.Conf.Mail.AttachmentFormat == "json" {
		if summaryFileName, err = template.CreateJsonFile(&summary); err != nil {
			log.Panic(err)
		}
	} else {
		if summaryFileName, err = template.CreateHTMLFile(&summary, config.Conf.Mail.TemplateDir); err != nil {
			log.Panic(err)
		}
	}
	defer func() {
		if err = os.Remove(summaryFileName); err != nil {
			log.Fatal(err)
		}
	}()

	// todo error handling
	sendMail("There are unexpected diffs in your system. See attachment for details!", summaryFileName, *config.Conf)
	log.Print("summary sent to " + config.Conf.Mail.ReceiverAdr)
}
