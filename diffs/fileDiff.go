package diffs

import (
	"bufio"
	"bytes"
	"goBackupInspector/backup"
	"goBackupInspector/summary"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const chunkSize = 64000

//CreateFileDiff unpacks zip file in <code>dumpDir directory</code> and compares the
// files between sysDir and the unpacked files
func CreateFileDiff(backupFilesDir, sysDir, pathInBackup, diffIgnoreFile string) (diffSummary summary.FileDiffSummary, err error){
	backupFilesDir = changeSuffix(backupFilesDir)
	sysDir = changeSuffix(sysDir)
	pathInBackup = changeSuffix(pathInBackup)
	if !strings.HasPrefix(pathInBackup, "/") {
		pathInBackup = "/" + pathInBackup
	}

	backupFile, err := backup.FindLastBackup(backupFilesDir)
	if err != nil {
		return
	}
	currentTime := time.Now()
	unpackDirectory := backupFilesDir + currentTime.Format("2006-01-02")
	defer func() {
		_ = backup.RemoveTMPDumpDir(unpackDirectory)
	}()
	// delete the temporary directory if it exists
	if _, errExists := os.Open(unpackDirectory); errExists == nil {
		if errExists = backup.RemoveTMPDumpDir(unpackDirectory); errExists != nil {
			return
		}
	}
	if err = backup.Unpack(backupFile, unpackDirectory); err != nil {
		return
	}

	diffSummary = fileDiff(unpackDirectory + pathInBackup, sysDir, diffIgnoreFile)
	diffSummary.BackupFileName = backupFile
	return
}

func changeSuffix(directory string) string{
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}
	return directory
}

func fileDiff(dumpDir, sysDir, diffIgnoreFile string) (diffSummary summary.FileDiffSummary){
	ignoreItems,_ := readDiffIgnore(diffIgnoreFile)
	ignoreRegex,_ := createIgnoreRegex(ignoreItems)

	dumpFiles, dirsInDump, ignoredElements := createFileDirChecklist(dumpDir, ignoreRegex)
	sysFiles, dirsInSys, ignoredElementsSys := createFileDirChecklist(sysDir, ignoreRegex)
	ignoredElements = append(ignoredElements, ignoredElementsSys...)

	// compare dumpFiles vs. sysFiles
	diffSummary = summary.FileDiffSummary{ LeftDir: dumpDir, RightDir: sysDir, FilesNotInDir: map[string][]string{},
											DirectoriesNotInDir: map[string][]string{},
											ComparedFiles: []summary.FileTuple{}, UnequalFiles: []summary.FileTuple{},
											IgnoredElement: ignoredElements}
	diffSummary.ComparedFiles, diffSummary.FilesNotInDir[sysDir] = findCuts(dumpDir, dumpFiles, sysDir, sysFiles)
	_, diffSummary.FilesNotInDir[dumpDir] = findCuts(sysDir, sysFiles, dumpDir, dumpFiles)
	_, diffSummary.DirectoriesNotInDir[sysDir] = findCuts(dumpDir, dirsInDump, sysDir, dirsInSys)
	_, diffSummary.DirectoriesNotInDir[dumpDir] = findCuts(sysDir, dirsInSys, dumpDir, dirsInDump)
	diffSummary.UnequalFiles = compareFiles(diffSummary.ComparedFiles)
	diffSummary.WithDifferences = diffSummary.HasDifferences()
	return
}

func compareFiles(filesToCompare []summary.FileTuple) (filesDiffs []summary.FileTuple) {
	for _, tuple := range filesToCompare {
		if !deepCompare(tuple.LeftFile, tuple.RightFile) {
			filesDiffs = append(filesDiffs, tuple)
		}
	}
	return
}

func deepCompare(file1, file2 string) bool {
	// src https://play.golang.org/p/YyYWuCRJXV

	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = f1.Close()
	}()

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = f2.Close()
	}()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

func findCuts(leftRootDir string, leftFilesNames []string, rightRootDir string, rightFileNames []string)  (filesToCompare []summary.FileTuple, filesNotInRightDir []string){
	for _, leftFile := range leftFilesNames {
		if rightFile := findFileWithName(leftFile, leftRootDir, rightFileNames, rightRootDir); rightFile != "" {
			filesToCompare = append(filesToCompare, summary.FileTuple{LeftFile: leftFile, RightFile: rightFile})
		} else {
			filesNotInRightDir = append(filesNotInRightDir, leftFile)
		}
	}
	return
}

func findFileWithName(leftFileName string, leftRootDir string, rightFileNames []string, rightRootDir string) (file string) {
	for _, file := range rightFileNames {
		if strings.TrimPrefix(file, rightRootDir) == strings.TrimPrefix(leftFileName, leftRootDir) {
			return file
		}
	}
	return ""
}

func createIgnoreRegex(ignoreItems []string) (ignoreRegex []*regexp.Regexp, err error) {
	for _, ignoreItem := range ignoreItems {
		regex, err := regexp.Compile(ignoreItem)
		if err != nil {
			return []*regexp.Regexp{}, err
		}
		ignoreRegex = append(ignoreRegex, regex)
	}
	return
}

func readDiffIgnore(diffIgnoreFile string) (ignoreItems []string, err error) {
	var file *os.File
	file, err = os.Open(diffIgnoreFile)
	if err != nil {
		return
	}
	defer func() { err = file.Close() }()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		ignoreItems = append(ignoreItems, scanner.Text())
	}

	err = scanner.Err()
	return
}

func isIgnored(name string, ignoreRegex []*regexp.Regexp) (bool, string) {
	for _, ignoreItem := range ignoreRegex {
		if ignoreItem.MatchString(name) {
			return true, ignoreItem.String()
		}
	}
	return false, ""
}

func createFileDirChecklist(root string, ignoreRegex []*regexp.Regexp) (files, directories []string,
																		ignoredElements []summary.IgnoredElement) {
	ignoredElements = make([]summary.IgnoredElement,0)
	// filepath.Walk
	filesAll, directoriesAll, err := filePathWalkDir(root)
	if err != nil {
		panic(err)
	}
	for _, file := range filesAll {
		if ignore, ruleName := isIgnored(file, ignoreRegex); ignore {
			ignoredElements = append(ignoredElements, summary.IgnoredElement{IgnoredElement: file, CausedRule: ruleName})
		} else {
			files = append(files, file)
		}
	}

	for _, dir := range directoriesAll {
		if ignore, ruleName := isIgnored(dir, ignoreRegex); ignore {
			ignoredElements = append(ignoredElements, summary.IgnoredElement{IgnoredElement: dir, CausedRule: ruleName})
		} else if ignoreWithSlash, ruleNameSlash := isIgnored(dir + "/", ignoreRegex); ignoreWithSlash {
			ignoredElements = append(ignoredElements, summary.IgnoredElement{IgnoredElement: dir, CausedRule: ruleNameSlash})
		} else {
			directories = append(directories, dir)
		}
	}
	return
}

func filePathWalkDir(root string) (files []string, directories []string,err error) {
	//var files []string
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		} else {
			directories = append(directories, path)
		}
		return nil
	})
	return files, directories, err
}
