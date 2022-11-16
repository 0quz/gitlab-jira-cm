package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/0quz/gitlab-jira-cm/pkg/model"
)

// Struct for incoming request
type MergeRequest struct {
	Project          Project          `json:"project" validate:"required"`
	ObjectAttributes ObjectAttributes `json:"object_attributes" validate:"required"`
}

type Project struct {
	Id        int    `json:"id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Namespace string `json:"namespace" validate:"required"`
}

type ObjectAttributes struct {
	LastCommit   LastCommit `json:"last_commit"`
	Url          string     `json:"url" validate:"required"`
	Iid          int        `json:"iid" validate:"required"`
	SourceBranch string     `json:"source_branch" validate:"required"`
	TargetBranch string     `json:"target_branch" validate:"required,target_branch"`
	State        string     `json:"state" validate:"required,state"`
	Action       string     `json:"action" validate:"required"`
	Title        string     `json:"title" validate:"required"`
}

type LastCommit struct {
	Author Author `json:"author"`
}

type Author struct {
	Email string `json:"email"`
}

// Struct for Kafka feed
type ProduceGitlabKafka struct {
	User         KUser         `json:"user"`
	Project      KProject      `json:"project"`
	MergeRequest KMergeRequest `json:"merge_request"`
	Title        string        `json:"title"`
	Watchers     []string      `json:"watchers"`
	Assignee     string        `json:"assignee"`
}

type KUser struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type KProject struct {
	ProjectId   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Namespace   string `json:"namespace"`
}

type KMergeRequest struct {
	Iid string `json:"iid"`
	Url string `json:"url"`
}

// JSON to Struct converter for incoming request
func (mr *MergeRequest) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(mr)
}

// Add an event
func (ms *MarcoService) NewEventAdd(mr MergeRequest, topicName string) error {
	e := &model.Event{
		MrUrl: mr.ObjectAttributes.Url,
		Topic: topicName,
	}
	return ms.CreateEvent(e)
}

// Add an hotfix error
func (ms *MarcoService) NewHotFixError(mr MergeRequest) error {
	h := &model.HotFixError{
		MrUrl: mr.ObjectAttributes.Url,
	}
	return ms.CreateHotFixError(h)
}

// Prepare requested data
func (pgk *ProduceGitlabKafka) ProcessedData(mr MergeRequest) {
	username, name := GetUserByEmail(mr.ObjectAttributes.LastCommit.Author.Email)
	watchers, assignee := GetWatchersAndAssignee(strconv.Itoa(mr.Project.Id), username)
	pgk.User.Name = name
	pgk.User.Username = username
	pgk.Project.ProjectId = strconv.Itoa(mr.Project.Id)
	pgk.Project.ProjectName = mr.Project.Name
	pgk.Project.Namespace = mr.Project.Namespace
	pgk.MergeRequest.Iid = strconv.Itoa(mr.ObjectAttributes.Iid)
	pgk.MergeRequest.Url = mr.ObjectAttributes.Url
	pgk.Title = mr.ObjectAttributes.Title
	pgk.Watchers = watchers
	pgk.Assignee = assignee
}

// Struct to JSON converter for kafka
func ProduceGitlabKafkaToJson(mr MergeRequest) ([]byte, error) {
	pgk := ProduceGitlabKafka{}
	pgk.ProcessedData(mr)
	return json.Marshal(pgk)
}

// The HotFix merge request handler
// It approves mr immediately and feeds the kafka topic
func (ms *MarcoService) HotFix(mr MergeRequest) error {
	err := ms.NewEventAdd(mr, os.Getenv("topic_urgent"))
	if err != nil {
		fmt.Println("NewEventAdd: failed executing hot-fix db add: ", err.Error())
		return err
	}
	// sleep X seconds for preparing mr
	time.Sleep(5 * time.Second)
	url := os.Getenv("gitlab_api_url") + "projects/" +
		strconv.Itoa(mr.Project.Id) + "/merge_requests/" +
		strconv.Itoa(mr.ObjectAttributes.Iid) + "/approve"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("NewRequest: New POST request error: ", err.Error())
		return err
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("gitlab_token"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("DefaultClient: New POST response error: ", err.Error())
		return err
	}
	if res.StatusCode != 201 {
		ms.NewHotFixError(mr)
		fmt.Println("Hot-Fix was not able to be approved. Check the mr url: " + mr.ObjectAttributes.Url)
		return nil
	}
	p, err := ProduceGitlabKafkaToJson(mr)
	if err != nil {
		fmt.Println("ProduceGitlabKafkaToJson: hot-fix marshal error: ", err.Error())
		return err
	}
	err = ProduceKafkaTopic(os.Getenv("topic_urgent"), p)
	if err != nil {
		fmt.Println("ProduceKafkaTopic: Produce kafka topic: ", os.Getenv("topic_urgent"), " error: ", err.Error())
		return err
	}
	//fmt.Println(string(p))
	return nil
}

// The standard merge request handler
// It feeds the kafka topic
func (ms *MarcoService) Standard(mr MergeRequest) error {
	err := ms.NewEventAdd(mr, os.Getenv("topic_standard"))
	if err != nil {
		fmt.Println("NewEventAdd: failed executing standard db add: ", err.Error())
		return err
	}
	p, err := ProduceGitlabKafkaToJson(mr)
	if err != nil {
		fmt.Println("ProduceGitlabKafkaToJson: Standard marshal error: ", err.Error())
		return err
	}
	err = ProduceKafkaTopic(os.Getenv("topic_standard"), p)
	if err != nil {
		fmt.Println("ProduceKafkaTopic: Produce kafka topic: ", os.Getenv("topic_standard"), " error: ", err.Error())
		return err
	}
	//fmt.Println(string(p))
	return nil
}
