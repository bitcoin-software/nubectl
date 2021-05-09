// Mock ssh file until golang-based ssh terminal work is stabilized.
package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"io"
	"strings"
	"syscall"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

// Golang has no built-in realpath(1) (e.g. for tilde expansion)
func expand(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}

// pass SSH args to external ssh binary
func doSSH(sshKeyPath string, username string, hostname string, port int) {
	sshString := fmt.Sprintf("ssh -i %s %s@%s -p%d", sshKeyPath, username, hostname, port)
	fmt.Fprintln(os.Stderr,"[debug] using sshkey: ", sshKey)
	fmt.Fprintln(os.Stderr,"[debug] ssh str:", sshString)

	sshArgs := strings.Fields(sshString)

	//syscall.Exec req for full/realpath to binaries
	binary, lookErr := exec.LookPath("ssh")
	if lookErr != nil {
		panic(lookErr)
	}

	syscall.Exec(binary, sshArgs, os.Environ())
}


func sshResource(name string, keyID string) {

	var statusurl string

	sshKeyPath, _ := expand(*sshkey)

	if !fileExists(sshKeyPath) {
		fmt.Printf("no such SSHKEY env or --sshkey: %s\n", sshKeyPath)
		os.Exit(1)
	}

	//get status for 'name' environment
	statusurl = cloudUrl + "/api/v1/status/" + name

	request, error := http.NewRequest("GET", statusurl, nil)
	request.Header.Set("cid", keyID)
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	decoder := json.NewDecoder(strings.NewReader(string(body)))

	for {
		var m interface{}
		if err := decoder.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintln(os.Stderr,"fatal:")
			//log.Fatal(err)
		}
//		fmt.Fprintln(os.Stderr,"res:", m)

//		md, ok := m.(map[string]interface{})
		md, _ := m.(map[string]interface{})

		value, exists := md["instanceid"]
		// In case when key is not present in map variable exists will be false.
		if exists {
			//fmt.Printf("key exists in map: %t, value: %v \n", exists, value)
			if value == name {
				fmt.Printf("key exists in map: %t, value: %v \n", exists, value)
				ssh_user, _ := md["ssh_user"]
				ssh_host, _ := md["ssh_host"]
				ssh_port, _ := md["ssh_port"]

				//  need type assertion from interface
				username, _ := ssh_user.(string)
				hostname, _ := ssh_host.(string)
				ssh_port_int, _ := ssh_port.(float64)
				var port int = int(ssh_port_int)

				fmt.Printf("%s@%s -p %d\n", ssh_user,ssh_host, port)
				doSSH(sshKeyPath, username, hostname, port)
			}
		}
	}
}

// interactive build/select list via promptui.Select
// from getStatus -> instanceid name
func sshSelectResource(keyID string) {
	//wip
}
