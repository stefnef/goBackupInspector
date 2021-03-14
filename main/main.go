package main

import (
	"errors"
	"fmt"
	"goBackupInspector/config"
	"goBackupInspector/diffs"
	"goBackupInspector/notify"
	"goBackupInspector/template"
	"log"
	"os"
)

func main() {
	var summaryFileName string
	var logFile *os.File
	var err error

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

	if logFile, err = createLogFile(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer func() {
		_ = logFile.Close()
	}()

	summary, err := diffs.CreateFileDiff(config.Conf.BackupFilesDir, config.Conf.SysDir,
		config.Conf.RelPathInBackup, config.Conf.DiffIgnoreFile)
	if err != nil {
		log.Panic(err)
	}
	if !summary.WithDifferences || !notify.WithUserNotification(summary, config.Conf.SummaryDir) {
		log.Print("summary created but not sent")
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

	notify.SaveSummary(summary, config.Conf.SummaryDir)

	defer func() {
		if err = os.Remove(summaryFileName); err != nil {
			log.Fatal(err)
		}
	}()

	// todo error handling
	notify.SendMail("There are unexpected diffs in your system. See attachment for details!", summaryFileName, *config.Conf)
	log.Print("summary sent to " + config.Conf.Mail.ReceiverAdr)
}

func createLogFile() (f *os.File, err error) {
	const maxSize int64 = 2048
	f, err = os.OpenFile("backupDifferences.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if f == nil {
		return nil, errors.New("could not open log file")
	}
	stat, _ := f.Stat()
	if stat.Size() >= maxSize {
		_ = f.Truncate(0)
	}
	log.SetOutput(f)
	return
}
