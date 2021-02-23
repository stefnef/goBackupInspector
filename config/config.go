package config

import (
	"encoding/json"
	"errors"
	"os"
)

type MailConfig struct {
	FromAdr string `json:"FromAdr"`
	UserName string `json:"UserName"`
	Password string `json:"Password"`
	SMTPServer string `json:"SMTPServer"`
	SMTPPort int `json:"SMTPPort"`
	ReceiverAdr string `json:"ReceiverAdr"`
	ReceiverName string `json:"ReceiverName"`
}

type Config struct {
	BackupFilesDir      string     `json:"BackupFilesDir"`
	RelPathInBackup 	string 		`json:RelPathInBackup`
	SysDir         		string     `json:"SysDir"`
	DiffIgnoreFile 		string     `json:"DiffIgnoreFile"`
	Mail           		MailConfig `json:"Mail"`
}

var Conf *Config

func LoadConfiguration(file string) error {
	var c Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&c)
	if err == nil {
		Conf = &c
	}
	return err
}

func CheckConfig() error {
	if _,err := os.Open(Conf.SysDir); err != nil && os.IsNotExist(err) {
		return errors.New(`error in config file:
								The directory <sysDir> (` + Conf.SysDir +  `) does not exist`)
	}

	if _,err := os.Open(Conf.BackupFilesDir); err != nil && os.IsNotExist(err) {
		return errors.New(`error in config file: 
								The directory of your backup files does not exist
								(` + Conf.BackupFilesDir + `)`)
	}

	if _,err := os.Open(Conf.DiffIgnoreFile); err != nil && os.IsNotExist(err) {
		return errors.New(`error in config file: 
								The path of your ignore file is wrong
								(` + Conf.DiffIgnoreFile + `)`)
	}
	return nil
}