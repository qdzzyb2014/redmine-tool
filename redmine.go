package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const token = "40e43a15340b2338c94642c82ea16e92a17ff445"
const contentType = "application/json"
const baseURL = "https://decentfox.net:10443"
const test = "https://decentfox.net:10443/projects/tireo2o"
const testURL = "http://www.redmine.org"

type Info struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Issue struct {
	Id           int    `json:"id"`
	Status       Info   `json:"status"`
	Project      Info   `json:"project"`
	Tracker      Info   `json:"tracker"`
	Assigned_to  Info   `json:"assigned_to"`
	Priority     Info   `json:"priority"`
	Author       Info   `json:"author"`
	FixedVersion Info   `json:"fixed_version"`
	Parent       Info   `json:"parent"`
	Subject      string `json:"subject"`
	Description  string `json:"description"`
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

type IssuesResp struct {
	Issues []Issue
}

func (resp IssuesResp) String() string {
	rv := ""
	for _, i := range resp.Issues {
		rv += i.String() + "\n"
	}
	return rv + "\n"
}

type Project struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
	Description string `json:"description"`
}

func (p Project) String() string {
	return fmt.Sprintf("id: %v\tname: %v", p.Id, p.Name)
}

type ProjectResp struct {
	Projects []Project
}

func (resp ProjectResp) String() string {
	rv := ""
	for _, i := range resp.Projects {
		rv += i.String() + "\n"
	}
	return rv + "\n"
}

func main() {
	command := os.Args[1]
	switch command {
	case "issue":
		getIssusInfo()
	case "issues":
		getIssues()
	case "project":
		Progects()
	default:
		fmt.Println("wrong command")
	}
}

func getClinet() *http.Client {
	return &http.Client{}
}

func NewRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("ContentType", contentType)
	req.Header.Add("X-Redmine-API-Key", token)
	return req
}

func getIssues() {
	projectId := flag.Int("project", 0, "project id. run 'redmine-tool projects' learn more info")
	flag.CommandLine.Parse(os.Args[2:])
	url := baseURL + "/" + "issues.json" 
	if *projectId != 0 {
		url += "?" + "project_id=" + strconv.Itoa(*projectId)
	}
	req := NewRequest("GET", url)
	resp, err := getClinet().Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	issues := IssuesResp{}
	json.Unmarshal(body, &issues)
	fmt.Println(issues)
}


func getIssusInfo() {
	id := flag.Int("id", 0, "")
	flag.CommandLine.Parse(os.Args[2:])
	if *id == 0 {
		fmt.Printf("Wrong id: %d\n", id)
		return
	}
	url := baseURL + "/" + "issues" + "/" + strconv.Itoa(*id) + ".json"
	req := NewRequest("GET", url)
	resp, err := getClinet().Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	issue := IssueResp{}
	json.Unmarshal(body, &issue)
	clipboard.WriteAll(issue.Is.GitMessage())
	fmt.Println(issue.Is)
}

func Progects() {
	url := baseURL + "/" + "projects.json"
	req := NewRequest("GET", url)
	resp, err := getClinet().Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	projects := ProjectResp{}
	json.Unmarshal(body, &projects)
	fmt.Println(projects)
}
