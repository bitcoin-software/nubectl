package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

// return true of dir/file exist
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func applyConfig(pubkey string) {
	envir := readCloudConfig()

	for nr, vm := range envir.Vm {
		fmt.Println(nr)
		fmt.Println("creating " + vm.Name)
		cpucount, _ := strconv.Atoi(vm.Cpu)
		createVM(vm.Image, cpucount, vm.Ram, vm.Disksize, pubkey, vm.Name)
	}

	for nr, jail := range envir.Container {
		fmt.Println(nr)
		fmt.Println("creating " + jail.Name)

		v := reflect.ValueOf(jail)
		typeOfS := v.Type()

		for i := 0; i < v.NumField(); i++ {
			fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
		}

		createJail(jail.Disksize, pubkey, jail.Name)
	}
}

func divertConfig(token string) {
	envir := readCloudConfig()

	for nr, vm := range envir.Vm {
		fmt.Println(nr)
		fmt.Println("deleting " + vm.Name)
		destroyResource(vm.Name, token)
	}

	for nr, jail := range envir.Container {
		fmt.Println(nr)
		fmt.Println("deleting " + jail.Name)
		destroyResource(jail.Name, token)
	}

}

func main() {

	var keypath string

	if len(os.Args) < 2 {
		fmt.Println("no arguments supplied! run 'nubectl --help' to get list of args")
		os.Exit(1)
	}

	flag.Parse()

	// keypath: get from args
	if len(*cloud_key) > 1 {
		keypath = *cloud_key
		//fmt.Println("hello world")
	} else {
		// keyppath: get from env(1)
		keypath = os.Getenv("CLOUD_KEY")
	}

	if !fileExists(keypath) {
		fmt.Printf("no such CLOUD_KEY env or --cloud_key: %s\n", keypath)
		os.Exit(1)
	}

	// cloudUrl: get from args
	if len(*cloud_url) > 1 {
		cloudUrl = *cloud_url
	} else {
		// keyppath: get from env(1)
		cloudUrl = os.Getenv("CLOUD_URL")
	}

	// sshKey: get from args
	if len(*ssh_key) > 1 {
		sshKey = *ssh_key
	} else {
		// sshKey: get from env(1)
		sshKey = os.Getenv("SSHKEY")
	}

	if len(cloudUrl) < 2 {
		fmt.Printf("no such CLOUD_URL env or --cloud_url\n")
		os.Exit(1)
	}

	apitoken := getToken(keypath)

	pubkey := string(readKey(keypath))

	time.Sleep(1 * time.Second)
	command := os.Args[1]

	if command == "create" {
		resourceType := os.Args[2]
		if resourceType == "vm" {
			createVmDialogue(pubkey, os.Args[3])
		} else if resourceType == "container" {
			createJailDialogue(pubkey, os.Args[3])
		} else {
			fmt.Println("Usage: nubectl create [vm|container]")
		}

	} else if command == "status" {
		if len(os.Args) == 3 {
			getStatus(os.Args[2], apitoken)
		} else {
			// empty or wrong arg num: show cluster status only
			getStatus("", apitoken)
		}
	} else if command == "destroy" {
		destroyResource(os.Args[2], apitoken)
	} else if command == "list" {
		listCluster(apitoken)
	} else if command == "apply" {
		applyConfig(pubkey)
	} else if command == "divert" {
		divertConfig(apitoken)
	} else if command == "ssh" {
		if len(cloudUrl) < 2 {
			fmt.Printf("no such SSHKEY env or --ssh_key\n")
			os.Exit(1)
		}
		if len(os.Args) == 3 {
			sshResource(os.Args[2], apitoken)
		} else {
			sshSelectResource(apitoken)
		}
	}

	//fmt.Println(pubkey)
}
