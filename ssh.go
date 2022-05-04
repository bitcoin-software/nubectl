//
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	//not for win?
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

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
	var sshExternal bool
	sshExternal = true

	if sshExternal {
		sshString := fmt.Sprintf("ssh -i %s %s@%s -p%d", sshKeyPath, username, hostname, port)
		fmt.Fprintln(os.Stderr, "[debug] using ssh_key: ", sshKey)
		fmt.Fprintln(os.Stderr, "[debug] ssh str:", sshString)

		sshArgs := strings.Fields(sshString)

		//syscall.Exec req for full/realpath to binaries
		binary, lookErr := exec.LookPath("ssh")
		if lookErr != nil {
			panic(lookErr)
		}

		syscall.Exec(binary, sshArgs, os.Environ())
	} else {
		server := fmt.Sprintf("%s:%d", hostname, port)

		publicKey, err := PublicKeyFile(sshKeyPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// for Password method method:
		// ssh.Password(pass),
		// todo: pass-protected key?
		//   ssh: this private key is passphrase protected
		config := &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				publicKey,
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		client, err := ssh.Dial("tcp", server, config)
		if err != nil {
			panic("Failed to dial: " + err.Error())
		}
		defer client.Close()

		// Each ClientConn can support multiple interactive sessions,
		// represented by a Session.
		session, err := client.NewSession()
		if err != nil {
			panic("Failed to create session: " + err.Error())
		}
		defer session.Close()

		fd := int(os.Stdin.Fd())
		state, err := terminal.MakeRaw(fd)
		if err != nil {
			fmt.Println(err)
		}
		defer terminal.Restore(fd, state)

		w, h, err := terminal.GetSize(fd)
		if err != nil {
			fmt.Println(err)
		}

		modes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}

		err = session.RequestPty("xterm", h, w, modes)
		if err != nil {
			fmt.Println(err)
		}

		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		session.Stdin = os.Stdin

		err = session.Shell()
		if err != nil {
			fmt.Println(err)
		}

		signal_chan := make(chan os.Signal, 1)

		//not for win?
		signal.Notify(signal_chan, syscall.SIGWINCH)

		go func() {
			for {
				s := <-signal_chan
				switch s {
				//not for win
				case syscall.SIGWINCH:
					fd := int(os.Stdout.Fd())
					w, h, _ = terminal.GetSize(fd)
					session.WindowChange(h, w)
				}
			}
		}()

		err = session.Wait()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func sshResource(name string, keyID string) {

	var statusurl string

	sshKeyPath, _ := expand(*ssh_key)

	if !fileExists(sshKeyPath) {
		fmt.Printf("no such SSHKEY env or --ssh_key: %s\n", sshKeyPath)
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
			fmt.Fprintln(os.Stderr, "fatal:")
			//log.Fatal(err)
		}
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

				fmt.Printf("%s@%s -p %d\n", ssh_user, ssh_host, port)
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
