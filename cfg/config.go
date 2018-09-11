package cfg

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
)

func readTomlFile(fp string, v interface{}) error {
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		if pe, ok := err.(*os.PathError); ok && pe.Err == syscall.ENOENT {
			return os.ErrNotExist
		}
		return ErrFileOpen
	}
	_, err = toml.Decode(string(data), v)
	if err != nil {
		return ErrTomlParse
	}
	return nil
}

func createConfigFile(fp string, def interface{}) error {
	file, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	fmt.Println("An config file with default values has been created.\n" +
		"Please edit it then run this again.")
	return toml.NewEncoder(file).Encode(def)
}

// Give a def (default value) as a start point to parse.
//
// Just like encoding/json, the given def will be modified in place.
//
// The 1st argument will be used as the config file path.
// If no file can be found there,
// a new file will be created there with the given "def" as content.
func Get(argv []string, def interface{}) error {
	fp, err := getFirstArg(argv)
	if err != nil {
		fmt.Println("Error:", ErrEnoughArgs)
		printHelpMessage()
		return ErrEnoughArgs
	}
	if fp == "version" {
		printVersion()
		return ErrPrintVersion
	}

	//var absFp string
	//if path.IsAbs(fp) {
	//	absFp = fp
	//} else {
	//	pwd, _ := os.Getwd()
	//	absFp = path.Join(pwd, fp)
	//}

	if strings.HasSuffix(fp, ".toml") != true {
		fmt.Println("Error:", ErrNotToml)
		printHelpMessage()
		return ErrNotToml
	}

	err = readTomlFile(fp, def)
	if err == nil {
		return nil
	}

	if err != os.ErrNotExist {
		fmt.Println("Error:", err)
		printHelpMessage()
		return err
	}

	err = createConfigFile(fp, def)
	if err == nil {
		return ErrExampleCreated
	}

	fmt.Println("Error:", err)
	printHelpMessage()
	return ErrFileCreate
}
