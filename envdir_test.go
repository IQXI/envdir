package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestEnvironmentFormatter(t *testing.T) {
	type args struct {
		env string
		cmd []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Command not found", args: args{env: "ok_dir", cmd: []string{"env"}}, wantErr: true},
		{name: "Env dir is file", args: args{env: "file.txt", cmd: []string{"ping"}}, wantErr: true},
		{name: "Env dir is empty", args: args{env: "empty_dir", cmd: []string{"ping"}}, wantErr: true},
		{name: "Command is missed", args: args{env: "empty_dir", cmd: []string{}}, wantErr: true},
		{name: "Ping", args: args{env: "ok_dir", cmd: []string{"ping", "127.0.0.1"}}, wantErr: false},
		{name: "Exit status 1", args: args{env: "ok_dir", cmd: []string{"ping"}}, wantErr: true},
	}

	//не понятно как использовать темп файл и темп дир

	err := os.Mkdir("ok_dir", 0777)
	if err != nil {
		if _, err := os.Stat("ok_dir"); os.IsNotExist(err) {
			if !os.IsExist(err) {
				t.Errorf("Cannot create ok_dir directory %v", err)
			}
		}
	}
	err = ioutil.WriteFile("ok_dir\\env_1.txt", []byte("132456789"), 0644)
	err = ioutil.WriteFile("ok_dir\\env_2.txt", []byte("hello world"), 0644)
	defer func() {
		os.RemoveAll("ok_dir")
	}()

	err = os.Mkdir("empty_dir", 0777)
	if err != nil {
		if _, err := os.Stat("empty_dir"); os.IsNotExist(err) {
			if !os.IsExist(err) {
				t.Errorf("Cannot create empty_dir directory %v", err)
			}
		}
	}
	defer func() {
		os.RemoveAll("empty_dir")
	}()

	err = ioutil.WriteFile("file.txt", []byte("Some content"), 0644)
	defer func() { fmt.Println("Removing file.txt"); os.Remove("file.txt") }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EnvironmentFormatter(tt.args.env, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("EnvironmentFormatter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
