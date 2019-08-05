package main

import (
	"errors"
	"fmt"
	"github.com/spf13/pflag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var env string
var cmd []string

func init() {
	pflag.StringVarP(&env, "envdir", "e", "", "ENV files directory")
	pflag.StringSliceVarP(&cmd, "command", "c", []string{}, "Command name")
}

func EnvironmentFormatter(env string, cmd []string) (err error) {
	var cmdEnv []string

	//если len(cmd) ноль, значит не введена комманда
	if len(cmd) == 0 {
		err = errors.New("executable command is missed!")
		log.Println(err)
		return err
	}

	if env != "" {
		dir, err := os.Stat(env)
		if os.IsNotExist(err) {
			err = errors.New(fmt.Sprintf("Directory %v is not exist!", env))
			log.Println(err)
			return err
		}
		if !dir.IsDir() {
			err = errors.New(fmt.Sprintf("%v is not direcrtory!", env))
			log.Println(err)
			return err
		}

		files, err := ioutil.ReadDir(env)
		if err != nil {
			log.Println(err)
			return err
		}
		if len(files) == 0 {
			err = errors.New(fmt.Sprintf("Directory %v is empty!", env))
			log.Println(err)
			return err
		}

		for _, f := range files {
			filePath := path.Join(env, f.Name())
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				log.Println(err)
				continue
			}
			if !fileInfo.IsDir() {
				envValue, err := ioutil.ReadFile(filePath)
				if err != nil {
					log.Println(err)
					continue
				} else {
					value := strings.TrimSuffix(fileInfo.Name(), filepath.Ext(fileInfo.Name())) + "=" + string(envValue)
					cmdEnv = append(cmdEnv, value)
				}
			}
		}
	}

	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if len(cmdEnv) > 0 {
		command.Env = cmdEnv
	}
	err = command.Start()
	err = command.Wait()
	fmt.Printf("Command: %v exit with error: %v\n", cmd, err)
	return err
}
