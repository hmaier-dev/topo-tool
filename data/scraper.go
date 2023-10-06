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

var switchesList = []SwitchInfo{
	{"sw_a-nord", "192.168.132.125"},
	{"sw_a-sued", "192.168.132.126"},
	//{"sw_c-nord-core", "192.168.132.120"},
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

// SwitchData represents data from the ssh-output
type SwitchData struct {
	MacAddress string //`json:"mac-address"`
	Vlan       string //`json:"vlan"`
	State      string //`json:"state"`
	Interface  string //`json:"interface"`
	Stack      string //`json:"stack"`
	Port       string //`json:"port"`
	Aging      string //`json:"aging"`
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
		Hostname *string `xml:"hostname,attr"` // using a *ptr because hostname could be nil
	} `xml:"EndpointProfile"`
}

var dbConn *sql.DB
var host = "localhost"

//var host = "db"

func init() {
	fmt.Printf("[%v] Connecting to the database...\n", time.Now())
	dbConn = dbConnect()
}

func main() {
	for {
		// TODO: there is a more elegant way to format the time
		fmt.Printf("[%v] Querying clearpass...\n", time.Now().Format("2006-01-02 15:04:05"))
		queryClearpass() // get all hostnames and ip-addresses
		fmt.Printf("[%v] Querying access-switches...\n", time.Now().Format("2006-01-02 15:04:05"))
		for _, sw := range switchesList {
			fmt.Printf("[%v] Talking to %v \n", time.Now().Format("2006-01-02 15:04:05"), sw.Name)
			queryAccessSwitches(sw)
		}
		fmt.Printf("[%v] Wating 30 Minutes for the next cycle...\n", time.Now().Format("2006-01-02 15:04:05")) // I don't know why this format is used and not the standard
		time.Sleep(30 * time.Minute)
	}
}

// ---------------------------------------------------------------------------------
// DATABASE-CONNECTION HANDLING
func dbConnect() *sql.DB {
	var user = "www-data"
	var password = "password123"
	var port = "3306"
	var dbName = "topology-tool"
	var source = user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName
	timeout := 30
	timeoutSec := 30 * time.Second
	wait := 5
	maxi := timeout / wait
	try := 1
	startTime := time.Now()
	dbConn, err := sql.Open("mysql", source)
	if err != nil {
		log.Fatal("sql.open failed: ", err)
	}

	for {
		err = dbConn.Ping()
		if err != nil {
			fmt.Printf("[%d/%d] Ping failed...\n", try, maxi)
			try += 1
		}
		if err == nil {
			fmt.Printf("[%v] Connecting to database successful!\n", time.Now().Format("2006-01-02 15:04:05"))
			break
		}
		if time.Since(startTime) >= timeoutSec {
			log.Fatal("Timeout: No connection to db.")
		}
		time.Sleep(5 * time.Second)
	}
	return dbConn
}

// -----------------------------------------------------
//
//	SSH-CONNECTION HANDLING
func queryAccessSwitches(sw SwitchInfo) {
	username, password := readCredFromFile("access_switch", "./cred.txt")
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Establish connection
	sshConn, err := ssh.Dial("tcp", sw.IP+":22", config)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer func(conn *ssh.Client) {
		err := conn.Close()
		if err != nil {
			//log.Fatalf("Error closing connection: %v", err)
			fmt.Printf("Got little error, closing the connection: %v \n", err.Error()) // this error occures frequently but I don't know why... ??!? :(
		}
	}(sshConn)

	// New session
	session, err := sshConn.NewSession()
	if err != nil {
		log.Fatalf("Failed to establish connection: %v", err)
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {
			if err.Error() != "EOF" {
				log.Fatalf("Error when closing session: %v", err)
			}
		}
	}(session)

	raw := runCommand(session, "display mac-address")
	filtered := processResponse(raw) // sort out non-significant interfaces and break up data
	//json := convertToJson(sorted)  // serialize into []byte

	// Clearing the current table
	_, err = dbConn.Exec("TRUNCATE `" + sw.Name + "`;")
	if err != nil {
		log.Fatal("problems with query ", err)
	}

	// Searching for corresponding mac-addresses to find the hostname
	for _, entry := range filtered {
		mac := strings.ReplaceAll(entry.MacAddress, "-", "")
		resp := getAllColumns(dbConn, "SELECT * from clearpass WHERE mac= '"+mac+"';")
		if 0 < len(resp) {
			interfaceName := entry.Interface
			macAddress := resp[0][2]
			hostname := resp[0][1]
			vlan := entry.Vlan
			stack := entry.Stack
			interfaceNum := entry.Port
			query := fmt.Sprintf("INSERT INTO `%s` (interface_name,mac,hostname,vlan,stack,interface_num) VALUES (?,?,?,?,?,?);", sw.Name)
			_, err = dbConn.Exec(query, interfaceName, macAddress, hostname, vlan, stack, interfaceNum)
			if err != nil {
				log.Fatalf("Error while inserting: %v \n", err)
			}
		}

	}

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
				MacAddress: matches[1],
				Vlan:       matches[2],
				State:      matches[3],
				Interface:  matches[4],
				Stack:      stack,
				Port:       port,
				Aging:      matches[5],
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
	// Building the slice to unwrap in dbConn.Exec
	values := []interface{}{} // I have no clue why there 2x {}....
	var mac string
	var host string
	var placeholders string

	// Call the clearpass API
	request := setupRequest()
	body := sendRequest(request)

	//// Read from fetched xml file
	//body := readXmlFromFile()

	// Parsing
	var ApiResponse TipsApiResponse
	err = xml.Unmarshal(body, &ApiResponse)
	if err != nil {
		log.Fatalf("Error parsing xml: %v \n", err)
	}

	// Clearing the clearpass table
	_, err = dbConn.Exec("TRUNCATE clearpass;")
	if err != nil {
		log.Fatal("problems with query ", err)
	}

	for _, endpoint := range ApiResponse.Endpoints {
		// just care about hostnames -> no nil
		if endpoint.EndpointProfile.Hostname != nil {
			mac = endpoint.MacAddress
			host = *endpoint.EndpointProfile.Hostname
			placeholders += "(?, ?),"
			values = append(values, host, mac)
		}
	}
	placeholders = placeholders[:len(placeholders)-1] // delete the last comma
	query := fmt.Sprintf("INSERT INTO clearpass (hostname, mac) VALUES %s;", placeholders)

	// Inserting into database
	_, err = dbConn.Exec(query, values...)
	if err != nil {
		log.Fatal("problems with query ", err)
	}
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
// TODO: make this function more elegant
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
func getAllColumns(conn *sql.DB, query string) [][]string {
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

func readXmlFromFile() []byte {
	//Read XML content from file
	xmlFile, err := os.Open("./formated-output.xml")
	if err != nil {
		log.Fatalf("Error opening XML file: %v \n", err)
	}
	defer xmlFile.Close()

	body, err := ioutil.ReadAll(xmlFile) // need to find an alternative... cuz this deprecated
	if err != nil {
		log.Fatalf("Error reading XML file: %v \n", err)
	}
	return body
}
