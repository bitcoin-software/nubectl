package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"strconv"
)

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

	prompt = promptui.Select{
		Label: "Select VM image",
		Items: []string{"centos7", "centos8", "rocky8", "ubuntu", "debian", "freebsd_ufs", "freebsd_zfs", "openbsd", "netbsd"},
	}

	_, image, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	cpucount, _ := strconv.Atoi(cpus)
	createVM(image, cpucount, ram, disksize, pubkey, name)
}
