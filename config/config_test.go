package config

import (
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	if err := LoadConfiguration("../test/configTestFile"); err != nil {
		t.Error(err)
	}

	if Conf == nil {
		t.Error("Config file not loaded")
	}

	if err := LoadConfiguration("../test/TestFile"); err == nil {
		t.Error("No error when file does not exist")
	}

}