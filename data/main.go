package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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

type SwitchData struct {
	Mac       string `json:"mac-address"`
	Vlan      string `json:"vlan"`
	State     string `json:"state"`
	Interface string `json:"interface"`
	Aging     string `json:"aging"`
}

func main() {
	fmt.Println("Starting ssh connect...")
	ssh_connector()
}

// -----------------------------------------------------
// SSH-Connector
func ssh_connector() {

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
	raw := run_command(session, "display mac-address")
	process_response(raw)
}

func run_command(session *ssh.Session, cmd string) string {
	out, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Fatalf("Failed to establish a session: %v", err)
	}
	return string(out)
}
func process_response(dirty string) {

	var macTable []SwitchData
	// 								 MAC			  VLAN    State   Port    Aging
	re := regexp.MustCompile(`([0-9a-fA-F-]{14})\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)`) //match the mac-address and the strings and whitespaces afterward

	lines := strings.Split(dirty, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if matches := re.FindStringSubmatch(line); matches != nil {

			fmt.Printf("%v", matches)
			data := SwitchData{
				Mac:       matches[1],
				Vlan:      matches[2],
				State:     matches[3],
				Interface: matches[4],
				Aging:     matches[5],
			}
			macTable = append(macTable, data)
		}

	}
	fmt.Printf("%v", macTable)
	//macTable = append(macTable, data)

	//json, err := json2.Marshal(macTable)
	//if err != nil {
	//	fmt.Println("Error while converting byte to json: ", err)
	//}
	//fmt.Println(json)
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

// -----------------------------------------------------
// Handy tools
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
