package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

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

func createJail(disksize string, key string, name string) {

	cloudurl := os.Getenv("CLOUDURL")

	httpposturl := cloudurl + "/api/v1/create/" + name
	//fmt.Println("HTTP JSON POST URL:", httpposturl)

	var jsonData = []byte(`{
		"type": "jail",
  		"imgsize": "` + disksize + `",
  		"pubkey": "` + key + `"
	}`)

	fmt.Println(string(jsonData))

	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	//fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
}

func createVM(image string, cores int, ramsize string, disksize string, key string, name string) {
	cloudurl := os.Getenv("CLOUDURL")

	httpposturl := cloudurl + "/api/v1/create/" + name
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	var jsonData = []byte(`{
	  "type": "bhyve",
	  "imgsize": "` + disksize + `",
	  "ram": "` + ramsize + `",
	  "cpus": "` + strconv.Itoa(cores) + `",
	  "img": "` + image + `",
	  "pubkey": "` + key + `"
	}`)

	//fmt.Println(string(jsonData))

	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	//fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
}

func getStatus(name string, keyID string) {
	cloudurl := os.Getenv("CLOUDURL")

	statusurl := cloudurl + "/api/v1/status/" + name
	//fmt.Println("HTTP JSON POST URL:", statusurl)

	request, error := http.NewRequest("GET", statusurl, nil)
	request.Header.Set("cid", keyID)
	//fmt.Println("cid:", keyID)

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	//fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
}

func listCluster(keyID string) {
	cloudurl := os.Getenv("CLOUDURL")
	statusurl := cloudurl + "/api/v1/cluster"
	//fmt.Println("HTTP JSON POST URL:", statusurl)

	request, error := http.NewRequest("GET", statusurl, nil)
	request.Header.Set("cid", keyID)
	//fmt.Println("cid:", keyID)

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	//fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
}

func destroyResource(name string, keyID string) {

	cloudurl := os.Getenv("CLOUDURL")

	destroyurl := cloudurl + "/api/v1/destroy/" + name
	//fmt.Println("HTTP JSON POST URL:", destroyurl)

	request, error := http.NewRequest("GET", destroyurl, nil)
	request.Header.Set("cid", keyID)
	//fmt.Println("cid:", keyID)

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	//fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
}

func createJailDialogue(pubkey string) {
	prompt := promptui.Select{
		Label: "Select disk size",
		Items: []string{"10g", "20g", "40g", "80g", "160g"},
	}

	_, disksize, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	createJail(disksize, pubkey, "testJ")
}

func createVmDialogue(pubkey string) {
	prompt := promptui.Select{
		Label: "Select how many CPUs you need",
		Items: []int{1, 2, 4, 8},
	}

	_, cpus, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	prompt = promptui.Select{
		Label: "Select how much RAM you need",
		Items: []string{"512m", "2g", "4g", "8g"},
	}

	_, ram, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	prompt = promptui.Select{
		Label: "Select disk size",
		Items: []string{"20g", "60g", "180g", "300g"},
	}

	_, disksize, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	prompt = promptui.Select{
		Label: "Select VM image",
		Items: []string{"centos7", "ubuntu20", "freebsd13", "docker"},
	}

	_, image, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf(cpus)
	createVM(image, 1, ram, disksize, pubkey, "testV")
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("no arguments supplied! run 'nubectl help' to get list of args")
		os.Exit(1)
	}

	//fmt.Println("hello world")
	keypath := os.Getenv("CLOUDKEY")

	apitoken := getToken(keypath)
	//fmt.Println(apitoken)
	//fmt.Println(`test ` + apitoken )
	//fmt.Println(string(readKey(keypath)))

	pubkey := string(readKey(keypath))
	/*
		createVM(1, "2g", "10g", pubkey, "testvm")
		time.Sleep(45 * time.Second)
	*/
	time.Sleep(1 * time.Second)
	command := os.Args[1]

	if command == "create" {
		resourceType := os.Args[2]
		if resourceType == "vm" {
			createVmDialogue(pubkey)
		} else if resourceType == "container" {
			createJailDialogue(pubkey)
		} else {
			fmt.Println("Usage: nubectl create [vm|container]")
		}

	} else if command == "status" {
		getStatus(os.Args[2], apitoken)
	} else if command == "destroy" {
		destroyResource(os.Args[2], apitoken)
	} else if command == "list" {
		listCluster(apitoken)
	}

	//fmt.Println(pubkey)
}
