package main

import (
	"bytes"
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
	assigned_to := fmt.Sprintf("Assigned To:%v %v", issue.Assigned_to.Name, issue.Assigned_to.Id)
	return fmt.Sprintf("#%v, %v\t%v\ndescription: \n%s", issue.Id, issue.Subject, assigned_to, issue.Description)
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
	case "working":
		updateIssue("working")
	default:
		fmt.Println("wrong command")
	}
}

func getClinet() *http.Client {
	return &http.Client{}
}

func Body(method, url string) []byte {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("ContentType", contentType)
	req.Header.Add("X-Redmine-API-Key", token)
	resp, err := getClinet().Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func getIssues() {
	projectId := flag.Int("project", 0, "project id. run 'redmine-tool projects' learn more info")
	flag.CommandLine.Parse(os.Args[2:])
	url := baseURL + "/" + "issues.json"
	if *projectId != 0 {
		url += "?" + "project_id=" + strconv.Itoa(*projectId)
	}
	issues := IssuesResp{}
	json.Unmarshal(Body("GET", url), &issues)
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
	issue := IssueResp{}
	json.Unmarshal(Body("GET", url), &issue)
	clipboard.WriteAll(issue.Is.GitMessage())
	fmt.Println(issue.Is)
}

func Progects() {
	url := baseURL + "/" + "projects.json"
	projects := ProjectResp{}
	json.Unmarshal(Body("GET", url), &projects)
	fmt.Println(projects)
}

func updateIssue(status string) {
	id := flag.Int("id", 0, "")
	flag.CommandLine.Parse(os.Args[2:])
	if *id == 0 {
		fmt.Printf("Wrong id: %d\n", id)
		return
	}
	url := baseURL + "/" + "issues" + "/" + strconv.Itoa(*id) + ".json"
	fmt.Println(url)
	issue := UpdateIssue{2, 14}
	fmt.Println(issue)
	fawData := UpdateIssueBody{issue}
	data, err := json.Marshal(fawData)
	fmt.Println(string(data))
	if err != nil {
		fmt.Println(err)
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("X-Redmine-API-Key", token)
	resp, err := getClinet().Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

type UpdateIssue struct {
	StatusId     int `json:"status_id"`
	AssignedToId int `json:"assigned_to_id"`
}

type UpdateIssueBody struct {
	Issue UpdateIssue `json:"issue"`
}
