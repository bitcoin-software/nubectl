package main

import (
	"crypto/md5"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

type Env struct {
	Version   string
	Container []Container
	Vm        []Vm
}

type Container struct {
	Name     string
	Disksize string
}

type Vm struct {
	Name     string
	Cpu      string
	Ram      string
	Disksize string
	Image    string
}

func readKey(keypath string) []byte {
	mykey, err := ioutil.ReadFile(keypath) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	keystr := string(mykey)
	//fmt.Printf(keystr)
	mykey = []byte(strings.TrimSuffix(keystr, "\n"))
	return mykey
}

func getToken(keypath string) string {

	mykey := readKey(keypath)
	result := fmt.Sprintf("%x", md5.Sum(mykey))
	return result
}

func readCloudConfig() Env {
	environment := Env{}

	yamlFile, err := ioutil.ReadFile("cloud.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &environment)

	fmt.Printf("%+v\n", environment)
	return environment
}
