package notify

import (
	"encoding/json"
	"goBackupInspector/summary"
	"io/ioutil"
	"log"
	"os"
)

const summaryFileName = "lastSummary.json"

func WithUserNotification(sum summary.FileDiffSummary, sumPath string) bool {
	f, err := os.Open(sumPath + "/" + summaryFileName)
	if err != nil && os.IsNotExist(err) {
		return true
	}
	defer func() { _ = f.Close() }()

	var other *summary.FileDiffSummary
	other, err = readSummary(sumPath)
	if err != nil {
		log.Print("error ", err)
		return true
	}

	return sum.Compare(*other) != int(summary.Equal)
}

func readSummary(sumPath string) (sum *summary.FileDiffSummary, err error) {
	var f *os.File
	var byteValue []byte
	sum = &summary.FileDiffSummary{}
	f, err = os.Open(sumPath + "/" + summaryFileName)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	if byteValue, err = ioutil.ReadAll(f); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(byteValue, sum); err != nil {
		return nil, err
	}
	return
}

func SaveSummary(sum summary.FileDiffSummary, sumPath string) {
	file, _ := json.MarshalIndent(sum, "", "\t")
	_ = ioutil.WriteFile(sumPath+"/"+summaryFileName, file, 0644)
}
