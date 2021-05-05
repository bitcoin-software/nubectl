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

func createVM(cores int, ramsize string, disksize string, key string, name string) {
	cloudurl := os.Getenv("CLOUDURL")

	httpposturl := cloudurl + "/api/v1/create/" + name
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	var jsonData = []byte(`{
	  "type": "bhyve",
	  "imgsize": "` + disksize + `",
	  "ram": "` + ramsize + `",
	  "cpus": "` + strconv.Itoa(cores) + `",
	  "img": "centos7",
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
		createJail("10g", pubkey, "testjail")
	} else if command == "status" {
		getStatus("testjail", apitoken)
	} else if command == "destroy" {
		destroyResource("testjail", apitoken)
	} else if command == "list" {
		listCluster(apitoken)
	}

	//fmt.Println(pubkey)
}
