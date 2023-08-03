package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
)

func main() {
	var username, password string
	username, password = get_cred("./cred.txt")
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

func get_cred(path string) (string, string) {
	var username, password string
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error when opening file: ", err)
		panic("error reading file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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
	return username, password
}
