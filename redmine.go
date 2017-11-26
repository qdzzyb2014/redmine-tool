package main

import (
	"strconv"
	"os"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"flag"
	"github.com/atotto/clipboard"
)

const token = "40e43a15340b2338c94642c82ea16e92a17ff445"
const contentType = "application/json"
const baseURL = "https://decentfox.net:10443"
const test = "https://decentfox.net:10443/projects/tireo2o"
const testURL = "http://www.redmine.org"

type Info struct {
	Id   int		`json:"id"`
	Name string		`json:"name"`
}

type Issue struct {
	Id            int 			`json:"id"`
	Status        Info			`json:"status"`
	Project       Info			`json:"project"`
	Tracker       Info			`json:"tracker"`
	Assigned_to   Info			`json:"assigned_to"`
	Priority      Info			`json:"priority"`
	Author        Info			`json:"author"`
	FixedVersion  Info			`json:"fixed_version"`
	Parent        Info			`json:"parent"`
	Subject       string		`json:"subject"`
	Description   string		`json:"description"`
}

func (issue Issue) String() string {
	return fmt.Sprintf("#%v, %v\ndescription: %s", issue.Id, issue.Subject, issue.Description)
}

func (issue Issue) GitMessage() string {
	return fmt.Sprintf("fixed #%v, %v", issue.Id, issue.Subject)
}

type IssueResp struct {
	Is Issue `json:"issue"`
}

func main() {
	command := os.Args[1]
	switch command {
	case "issues":
		getIssusInfo()
	default:
		fmt.Println("wrong command")
	}
}

func getClinet() *http.Client {
	// proxy := func(_ *http.Request) (*url.URL, error) {
	// 	return url.Parse("http://192.168.0.106:8001")
	// }
	// transport := &http.Transport{Proxy: proxy}
	return &http.Client{}
}

func getIssusInfo() {
	id := flag.Int("id", 0, "")
	flag.CommandLine.Parse(os.Args[2:])
	if *id == 0 {
		fmt.Printf("Wrong id: %d\n", id)
		return
	}
	url := baseURL + "/" + "issues" + "/" + strconv.Itoa(*id) + ".json"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("ContentType", contentType)
	req.Header.Add("X-Redmine-API-Key", token)
	resp, err := getClinet().Do(req)
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	issue := IssueResp{}
	json.Unmarshal(body, &issue)
	clipboard.WriteAll(issue.Is.GitMessage())
	fmt.Println(issue.Is)
}
