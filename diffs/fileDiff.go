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
func CreateFileDiff(backupFilesDir, sysDir, pathInBackup, diffIgnoreFile string) (diffSummary summary.FileDiffSummary, err error) {
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

	diffSummary = fileDiff(unpackDirectory+pathInBackup, sysDir, diffIgnoreFile)
	diffSummary.BackupFileName = backupFile
	return
}

func changeSuffix(directory string) string {
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}
	return directory
}

func fileDiff(dumpDir, sysDir, diffIgnoreFile string) (diffSummary summary.FileDiffSummary) {
	var filesToCompare []summary.FileTuple
	ignoreItems, _ := readDiffIgnore(diffIgnoreFile)
	ignoreRegex, _ := createIgnoreRegex(ignoreItems)

	dumpFiles, dirsInDump, ignoredElements := createFileDirChecklist(dumpDir, ignoreRegex)
	sysFiles, dirsInSys, ignoredElementsSys := createFileDirChecklist(sysDir, ignoreRegex)
	ignoredElements = append(ignoredElements, ignoredElementsSys...)

	// compare dumpFiles vs. sysFiles
	diffSummary = summary.FileDiffSummary{BackupDir: dumpDir, SystemDir: sysDir,
		Date: time.Now(),
		FilesNotInDir: map[string][]string{
			summary.DirBackup: make([]string, 0),
			summary.DirSystem: make([]string, 0),
		},
		DirectoriesNotInDir: map[string][]string{
			summary.DirBackup: make([]string, 0),
			summary.DirSystem: make([]string, 0),
		},
		ComparedFiles: []string{}, UnequalFiles: []summary.FileTuple{},
		IgnoredElements: ignoredElements}
	filesToCompare, diffSummary.FilesNotInDir[summary.DirSystem] = findCuts(dumpDir, dumpFiles, sysDir, sysFiles)
	_, diffSummary.FilesNotInDir[summary.DirBackup] = findCuts(sysDir, sysFiles, dumpDir, dumpFiles)
	_, diffSummary.DirectoriesNotInDir[summary.DirSystem] = findCuts(dumpDir, dirsInDump, sysDir, dirsInSys)
	_, diffSummary.DirectoriesNotInDir[summary.DirBackup] = findCuts(sysDir, dirsInSys, dumpDir, dirsInDump)
	diffSummary.UnequalFiles, diffSummary.ComparedFiles = compareFiles(filesToCompare)
	diffSummary.WithDifferences = diffSummary.HasDifferences()
	deletePrefixes(&diffSummary)
	return
}

//deletePrefixes Deletes backup and system paths from file names
func deletePrefixes(diff *summary.FileDiffSummary) {
	// files not in backup/system
	for idx, element := range diff.FilesNotInDir[summary.DirBackup] {
		diff.FilesNotInDir[summary.DirBackup][idx] = strings.TrimPrefix(element, diff.SystemDir)
	}
	for idx, element := range diff.FilesNotInDir[summary.DirSystem] {
		diff.FilesNotInDir[summary.DirSystem][idx] = strings.TrimPrefix(element, diff.BackupDir)
	}

	// directories not in backup/system
	for idx, element := range diff.DirectoriesNotInDir[summary.DirBackup] {
		diff.DirectoriesNotInDir[summary.DirBackup][idx] = strings.TrimPrefix(element, diff.SystemDir)
	}
	for idx, element := range diff.DirectoriesNotInDir[summary.DirSystem] {
		diff.DirectoriesNotInDir[summary.DirSystem][idx] = strings.TrimPrefix(element, diff.BackupDir)
	}

	// compared files
	for idx, element := range diff.ComparedFiles {
		noPrefix := strings.TrimPrefix(element, diff.BackupDir)
		noPrefix = strings.TrimPrefix(noPrefix, diff.SystemDir)
		diff.ComparedFiles[idx] = noPrefix
	}

	// ignored elements
	for idx, element := range diff.IgnoredElements {
		if strings.HasPrefix(element.IgnoredElement, diff.BackupDir) {
			element.IgnoredElement = summary.DirBackup + ": " + strings.TrimPrefix(element.IgnoredElement, diff.BackupDir)
		} else {
			element.IgnoredElement = summary.DirSystem + ": " + strings.TrimPrefix(element.IgnoredElement, diff.SystemDir)
		}
		diff.IgnoredElements[idx].IgnoredElement = element.IgnoredElement
	}

	// unequal files
	for idx, element := range diff.UnequalFiles {
		diff.UnequalFiles[idx].BackupFile = strings.TrimPrefix(element.BackupFile, diff.BackupDir)
		diff.UnequalFiles[idx].SystemFile = strings.TrimPrefix(element.SystemFile, diff.SystemDir)
	}
}

func compareFiles(filesToCompare []summary.FileTuple) (filesDiffs []summary.FileTuple, comparedFiles []string) {
	for _, tuple := range filesToCompare {
		comparedFiles = append(comparedFiles, tuple.BackupFile)
		if !deepCompare(tuple.BackupFile, tuple.SystemFile) {
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

func findCuts(leftRootDir string, leftFilesNames []string, rightRootDir string, rightFileNames []string) (filesToCompare []summary.FileTuple, filesNotInRightDir []string) {
	filesToCompare = make([]summary.FileTuple, 0)
	filesNotInRightDir = make([]string, 0)
	for _, leftFile := range leftFilesNames {
		if rightFile := findFileWithName(leftFile, leftRootDir, rightFileNames, rightRootDir); rightFile != "" {
			filesToCompare = append(filesToCompare, summary.FileTuple{BackupFile: leftFile, SystemFile: rightFile})
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
	ignoredElements = make([]summary.IgnoredElement, 0)
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
		} else if ignoreWithSlash, ruleNameSlash := isIgnored(dir+"/", ignoreRegex); ignoreWithSlash {
			ignoredElements = append(ignoredElements, summary.IgnoredElement{IgnoredElement: dir, CausedRule: ruleNameSlash})
		} else {
			directories = append(directories, dir)
		}
	}
	return
}

func filePathWalkDir(root string) (files []string, directories []string, err error) {
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
