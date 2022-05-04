package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// image sample
// curl https://bitclouds.sh/images
//{
//  "images": [
//    "ubuntu", 
//    "bitcoind", 
//    "centos", 
//    "clightning", 
//    "bsdjail", 
//    "lnd", 
//    "freebsd", 
//    "debian", 
//    "freebsd-ufs", 
//    "netbsd", 
//    "openbsd"
//  ]
//}
type images struct {
	Images []string
}

func createJailDialogue(pubkey string, name string) {
	prompt := promptui.Select{
		Label: "Select disk size",
		Items: []string{"10g", "20g", "40g", "80g", "160g"},
	}

	_, disksize, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	createJail(disksize, pubkey, name)
}

func createVmDialogue(pubkey string, name string) {
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

	// Start dynamic image list by CLOUD_URL/images URL
	// get image list from cloudurl/images
	var imageurl string
	imageurl = cloudUrl + "/images"

	nubeClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, imageurl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "nubectl")

	res, getErr := nubeClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	images1 := images{}
	jsonErr := json.Unmarshal([]byte(body), &images1)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	// End of Start dnynamic image list by  CLOUD_URL/images URL

	prompt = promptui.Select{
		Label: "Select VM image",
//		Items: []string{"centos7", "centos8", "ubuntu", "debian", "freebsd_ufs", "freebsd_zfs", "openbsd", "netbsd"},
		Items: images1.Images,
	}

	_, image, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	cpucount, _ := strconv.Atoi(cpus)
	createVM(image, cpucount, ram, disksize, pubkey, name)
}
