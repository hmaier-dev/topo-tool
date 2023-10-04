package main

import (
	"bufio"
	"database/sql" // you need to use a driver with this
	"encoding/base64"
	json2 "encoding/json"
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // this is the sql driver
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// SwitchInfo is used for the list containing the access switches
type SwitchInfo struct {
	Name string
	IP   string
}

// SwitchData represents data from the ssh-output
type SwitchData struct {
	Mac       string `json:"mac-address"`
	Vlan      string `json:"vlan"`
	State     string `json:"state"`
	Interface string `json:"interface"`
	Stack     string `json:"stack"`
	Port      string `json:"port"`
	Aging     string `json:"aging"`
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

// TipsApiResponse -------------------------------------------------------------------------------------------
// Clearpass API Response
type TipsApiResponse struct {
	Endpoints []Endpoint `xml:"Endpoints>Endpoint"`
}

// Endpoint everything which is nested inside (not between) a tag, is an attribute short: attr
type Endpoint struct {
	MacAddress      string `xml:"macAddress,attr"`
	EndpointProfile struct {
		Hostname *string `xml:"hostname,attr"` // I really don't know why this is not an attr?!?!?!
	} `xml:"EndpointProfile"`
	//Hostname   *string `xml:"EndpointProfile>hostname,attr"`
}

var conn *sql.DB

func init() {
	fmt.Printf("[%v] Connecting to the database...\n", time.Now())
	conn = dbConnect()
}

func main() {
	for {
		fmt.Printf("[%v] Querying clearpass...\n", time.Now())
		queryClearpass() // get all hostnames and ip-addresses
		fmt.Printf("[%v] Querying access-switches...\n", time.Now())
		queryAccessSwitches() //
		time.Sleep(30 * time.Minute)
	}
}

// Establishing the connection to the
func dbConnect() *sql.DB {
	var user = "www-data"
	var password = "password123"
	var host = "localhost"
	//var host = "db"
	var port = "3306"
	var dbName = "topology-tool"
	var source = user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName
	timeout := 30
	timeoutSec := 30 * time.Second
	wait := 5
	maxi := timeout / wait
	try := 1
	startTime := time.Now()
	conn, err := sql.Open("mysql", source)
	if err != nil {
		log.Fatal("sql.open failed: ", err)
	}

	for {
		err = conn.Ping()
		if err != nil {
			fmt.Printf("[%d/%d] Ping failed...\n", try, maxi)
			try += 1
		}
		if err == nil {
			fmt.Printf("[%v] Connecting to database successful!\n", time.Now())
			break
		}
		if time.Since(startTime) >= timeoutSec {
			log.Fatal("Timeout: No connection to db.")
		}
		time.Sleep(5 * time.Second)
	}
	return conn
}

// -----------------------------------------------------
//
//	the connection to the access switches
func queryAccessSwitches() {
	username, password := readCredFromFile("access_switch", "./cred.txt")
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
	defer func(conn *ssh.Client) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Error closing connection...")
		}
	}(conn)

	// New session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Failed to establish connection: %v", err)
	}
	raw := runCommand(session, "display mac-address")
	sorted := processResponse(raw) // sort out non-significant interfaces and break up data
	json := convertToJson(sorted)  // serialize into []byte

	fmt.Printf("%v", string(json))

}

func runCommand(session *ssh.Session, cmd string) string {
	out, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Fatalf("Failed to establish a session: %v", err)
	}
	return string(out)
}
func processResponse(dirty string) []SwitchData {
	var macTable []SwitchData
	// Explaining the regex:
	// 								 MAC			  VLAN    State   Port    Aging
	re := regexp.MustCompile(`([0-9a-fA-F-]{14})\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)`) // match the mac-address and the strings/whitespaces afterward
	lines := strings.Split(dirty, "\n")                                            // building lines from
	// Going through the WHOLE ssh output!
	for _, line := range lines {
		if matches := re.FindStringSubmatch(line); matches != nil {
			if !validate(matches) {
				continue
			}
			stack, port := disassembleInterfaceString(matches[4])
			data := SwitchData{
				Mac:       matches[1],
				Vlan:      matches[2],
				State:     matches[3],
				Interface: matches[4],
				Stack:     stack,
				Port:      port,
				Aging:     matches[5],
			}
			macTable = append(macTable, data)
		} // END if
	} // END for
	return macTable
}

func convertToJson(macTable []SwitchData) []byte {
	// Build []byte with json
	json, err := json2.Marshal(macTable)
	if err != nil {
		fmt.Println("Error while converting byte to json: ", err)
	}
	return json
}

func validate(matches []string) bool {
	re := regexp.MustCompile(`BAGG`) // Bridge-Aggregation
	if match := re.FindString(matches[4]); match != "" {
		return false
	}
	return true
}

func disassembleInterfaceString(interfaceStr string) (string, string) {
	re := regexp.MustCompile(`\d+`)                                // digits in interface string
	if match := re.FindAllString(interfaceStr, -1); match != nil { // FindAllString with the -1 works, but I don't know why...
		return match[0], match[2]
	} else {
		log.Fatalf("Could not disassemble: %s\n", interfaceStr)
		return "", ""
	}
}

//-----------------------------------------------------

// --------------------------
// CLEARPASS-CONNECTOR
func queryClearpass() {
	var err error
	//request := setupRequest()
	//// Requesting
	//body := sendRequest(request)

	// Read XML content from file
	xmlFile, err := os.Open("./formated-output.xml")
	if err != nil {
		fmt.Println("Error opening XML file:", err)
		return
	}
	defer xmlFile.Close()

	body, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading XML file:", err)
		return
	}

	// Parsing
	var ApiResponse TipsApiResponse
	err = xml.Unmarshal(body, &ApiResponse)
	if err != nil {
		log.Fatalf("Error parsing xml: %v \n", err)
	}

	// Clearing the clearpass table
	_, err = conn.Exec("TRUNCATE clearpass;")
	if err != nil {
		log.Fatal("problems with query ", err)
	}

	query := "INSERT INTO clearpass (hostname, mac) VALUES "
	values := []interface{}{} // I have no clue why there 2x {}....
	var mac string
	var host string
	// Iterating
	for _, endpoint := range ApiResponse.Endpoints {
		mac = endpoint.MacAddress
		host = " "
		if endpoint.EndpointProfile.Hostname != nil {
			host = *endpoint.EndpointProfile.Hostname
		}
		//query += fmt.Sprintf("(%s, %s), ", host, mac)
		query += "(?, ?),"
		values = append(values, host, mac)
	}
	query = query[:len(query)-1] // delete the last comma
	//query += ";"

	fmt.Println(query)

	// Inserting into database
	_, err = conn.Exec(query, values...)
	if err != nil {
		log.Fatal("problems with query ", err)
	}

	//// Print the response body
	//fmt.Println("Response:", string(body))
	// write the whole body at once

	//// Writing the xml to file
	//err = ioutil.WriteFile("output.txt", body, 0644)
	//if err != nil {
	//	panic(err)
	//}

}

func setupRequest() *http.Request {
	var username, password string
	username, password = readCredFromFile("clearpass", "./cred.txt")
	url := "http://clearpass-cl.ipb-halle.de/tipsapi/config/read/Endpoint"
	authString := username + ":" + password
	encodedAuthString := base64.StdEncoding.EncodeToString([]byte(authString))
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Authorization", "Basic "+encodedAuthString)
	return request
}
func sendRequest(request *http.Request) []byte {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body) // from http.Response into []byte
	if err != nil {
		log.Fatalf("Error reading response body: %v\n", err)
	}
	return body
}

// -----------------------------------------------------
// Handy tools
func readCredFromFile(forConnector string, path string) (string, string) {
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

// Return all rows from a database query
func getQuery(conn *sql.DB, query string) [][]string {
	rows, err := conn.Query(query)
	if err != nil {
		log.Fatal("problems with query", err)
	}
	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("Cannot address variable rows...")
		return [][]string{}
	}

	rawResult := make([][]byte, len(cols)) // [row][values] -> e.g. row: [[value][value][value]]
	dest := make([]interface{}, len(cols)) // .Scan() needs []any as result type
	allRows := make([][]string, 0)
	for i := range cols {
		dest[i] = &rawResult[i] // mapping dest indices to byte slice
	}
	for rows.Next() {
		err := rows.Scan(dest...)
		if err != nil {
			log.Fatal("problems scanning the database", err)
		}
		singleRow := make([]string, len(cols))
		for i, raw := range rawResult {
			singleRow[i] = string(raw) // from byte to string
			//fmt.Printf("%v -> %v \n", i, singleRow)
		}
		allRows = append(allRows, singleRow)
	}
	return allRows
}
