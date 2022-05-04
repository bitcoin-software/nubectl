package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	cloud_key = flag.String("cloud_key", "", "Path to CLOUD_KEY ( or CLOUD_KEY environment )")
	cloud_url = flag.String("cloud_url", "", "Cloud URL ( or CLOUD_URL environment )")
	ssh_key   = flag.String("ssh_key", "~/.ssh/id_ed25519", "SSH keypath for ssh ( or SSHKEY environment )")
)

var cloudUrl string
var sshKey string

func createJail(disksize string, key string, name string) {

	httpposturl := cloudUrl + "/api/v1/create/" + name
	//fmt.Println("HTTP JSON POST URL:", httpposturl)

	var jsonData = []byte(`{
		"image": "jail",
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

	httpposturl := cloudUrl + "/api/v1/create/" + name
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	var jsonData = []byte(`{
		"imgsize": "` + disksize + `",
		"ram": "` + ramsize + `",
		"cpus": "` + strconv.Itoa(cores) + `",
		"image": "` + image + `",
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

// show environment (when name is set) or cluster status
func getStatus(name string, keyID string) {

	var statusurl string

	if len(name) > 1 {
		//get status for 'name' environment
		statusurl = cloudUrl + "/api/v1/status/" + name
	} else {
		// when 'name' is absent then get cluster status
		statusurl = cloudUrl + "/api/v1/cluster"
	}

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

	fmt.Fprintln(os.Stderr, "response Status:", response.Status)
	//fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Fprintln(os.Stderr, "response Body:")
	fmt.Println(string(body))
}

func listCluster(keyID string) {
	statusurl := cloudUrl + "/api/v1/cluster"
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

	destroyurl := cloudUrl + "/api/v1/destroy/" + name
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
