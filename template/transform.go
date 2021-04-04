package template

import (
	"encoding/json"
	"goBackupInspector/summary"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func CreateHTMLFile(sum *summary.FileDiffSummary, templatePath string) (fileName string, err error) {
	var tpl *template.Template
	var f *os.File
	//tpl = template.Must(template.ParseFiles("summary.tpl", "header.tpl", "footer.tpl", "notInDir.tpl"))
	tpl, err = template.New("new.gohtml").Funcs(template.FuncMap{
		"ShowFilesNotInLeftDir": func(innerSum summary.FileDiffSummary) []string {
			return innerSum.FilesNotInDir[summary.DirBackup]
		},
		"ShowFilesNotInRightDir": func(innerSum summary.FileDiffSummary) []string {
			return innerSum.FilesNotInDir[summary.DirSystem]
		},

		"ShowDirsNotInLeftDir": func(innerSum summary.FileDiffSummary) []string {
			return innerSum.DirectoriesNotInDir[summary.DirBackup]
		},
		"ShowDirsNotInRightDir": func(innerSum summary.FileDiffSummary) []string {
			return innerSum.DirectoriesNotInDir[summary.DirSystem]
		},
		"LeftWithoutPath": func(tuple summary.FileTuple) string {
			return strings.TrimPrefix(tuple.BackupFile, sum.BackupDir)
		},
	}).ParseFiles(templatePath+"summary.tpl", templatePath+"header.tpl",
		templatePath+"footer.tpl", templatePath+"notInDir.tpl")
	if err != nil {
		return
	}

	fileName = "summary_" + sum.Date.Format("2006_01_02") + ".html"

	if f, err = os.Create(fileName); err != nil {
		return
	}
	defer func() {
		err = f.Close()
	}()
	if err = tpl.ExecuteTemplate(f, "summary.tpl", *sum); err != nil {
		return
	}
	return
}

func CreateJsonFile(sum *summary.FileDiffSummary) (fileName string, err error) {
	fileName = "diffSummary_" + time.Now().Format("2006_01_02") + ".json"
	file, _ := json.MarshalIndent(sum, "", "\t")
	_ = ioutil.WriteFile(fileName, file, 0644)
	return
}
