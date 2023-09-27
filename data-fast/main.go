package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net/http"
	"os"
)

type SwitchInfo struct {
	Name string
	IP   string
}

var switchesList = []SwitchInfo{
	{"sw_a-nord", "192.168.132.125"},
	{"sw_a-sued", "192.168.132.126"},
	{"sw_c-nord-core", "192.168.132.120"},
	{"sw_c-nord", "192.168.132.121"},
	{"sw_c-sued", "192.168.132.122"},
	{"sw_b", "192.168.132.132"},
	{"sw_r", "192.168.132.141"},
	{"sw_r2", "192.168.132.142"},
	{"sw_d", "192.168.132.127"},
	{"sw_e1", "192.168.132.138"},
	{"sw_e2", "192.168.132.123"},
	{"sw_f", "192.168.132.202"},
	{"sw_g", "192.168.132.146"},
	{"sw_n1", "192.168.132.134"},
	// {"SW_N2", "192.168.132.144"}, currently not reachable
	{"sw_ws", "192.168.132.130"},
	{"sw_kms", "192.168.132.128"},
	// {"SW_P", "192.168.132.136"}, 1/0/28 set as uplink
	// {"SW_A121", "192.168.132.131"}, 1/0/24 set as uplink
}

func main() {
	ssh_connector()
}

// -----------------------------------------------------
// SSH-Connector
func ssh_connector() {

	//for _, sw := range switchesList {
	//	fmt.Printf("Name: %v, IP: %v \n", sw.Name, sw.IP)
	//}

	username, password := get_cred("access_switch", "./cred.txt")
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Establish connection
	conn, err := ssh.Dial("tcp", switchesList[1].IP+":22", config)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()
	// New session

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Failed to establish connection: %v", err)
	}
	run := run_command(session, "display mac-address")
	println(string(run))

}

func run_command(session *ssh.Session, cmd string) []byte {
	out, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Fatalf("Failed to establish connection: %v", err)
	}
	return out
}

//-----------------------------------------------------

// --------------------------
// CLEARPASS-CONNECTOR
func clearpass() {
	var username, password string
	username, password = get_cred("clearpass", "./cred.txt")
	url := "http://clearpass-cl.ipb-halle.de/tipsapi/config/read/Endpoint"
	authString := username + ":" + password
	encodedAuthString := base64.StdEncoding.EncodeToString([]byte(authString))
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("Authorization", "Basic "+encodedAuthString)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	var arr []byte
	n, err := response.Body.Read(arr)
	fmt.Printf("%v \n", n)
	fmt.Printf("%v \n", arr)

	defer response.Body.Close()

}

func get_cred(forConnector string, path string) (string, string) {
	var username, password string
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error when opening file: ", err)
		panic("error reading file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	switch forConnector {
	case "clearpass":
		for scanner.Scan() {
			line := scanner.Text()
			_, err = fmt.Sscanf(line, "api_user=%s", &username)
			if err == nil {
				continue
			}
			_, err = fmt.Sscanf(line, "api_password=%s", &password)
			if err == nil {
				continue
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("error reading file: ", err)
		}
	case "access_switch":
		for scanner.Scan() {
			line := scanner.Text()
			_, err = fmt.Sscanf(line, "sw_user=%s", &username)
			if err == nil {
				continue
			}
			_, err = fmt.Sscanf(line, "sw_password=%s", &password)
			if err == nil {
				continue
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("error reading file: ", err)
		}

	}
	return username, password
}

//---------------------------------------------------------------------
