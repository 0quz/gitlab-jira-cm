package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Struct for response body
type UserDetail struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	State       string `json:"state"`
	AccessLevel int    `json:"access_level"`
}

// It returns true if the value matches in value of the slice
func ContainsInSclice(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// If the value matches in value of the slice. It extracts the value in the slice and returns it back
func RemoveInSlice(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// It shuffles slice
func ShuffleSlice(slice []string) {
	if len(slice) > 1 {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(slice), func(i, j int) { slice[i], slice[j] = slice[j], slice[i] })
	}
}

// It returns excluded usernames
func GetExcludedUsernames() []string {
	excludedUsernames := []string{}
	page := 1
	for {
		req, err := http.NewRequest("GET", os.Getenv("gitlab_api_url")+"users?page="+
			strconv.Itoa(page)+"&per_page=100&admins=true&active=true", nil)
		if err != nil {
			fmt.Println("getExcludedUsernames req error: ", err.Error())
			break
		}
		req.Header.Add("Authorization", "Bearer "+os.Getenv("gitlab_token"))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("getExcludedUsernames res error: ", err.Error())
			break
		}
		if res.StatusCode != 200 {
			fmt.Println("getExcludedUsernames status code error: ", res.StatusCode)
			break
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		ud := []UserDetail{}
		err = json.Unmarshal(resBody, &ud)
		if err != nil {
			fmt.Println("getExcludedUsernames unmashal error: ", err.Error())
			break
		}
		if reflect.DeepEqual(ud, []UserDetail{}) {
			break
		}
		for _, v := range ud {
			excludedUsernames = append(excludedUsernames, v.Username)
		}
		page++
	}
	return excludedUsernames
}

// It returns external usernames
func GetExternalUsernames() []string {
	externalUsernames := []string{}
	page := 1
	for {
		req, err := http.NewRequest("GET", os.Getenv("gitlab_api_url")+"users?page="+
			strconv.Itoa(page)+"&per_page=100&external=true", nil)
		if err != nil {
			fmt.Println("getExternalUsernames req error: ", err.Error())
			break
		}
		req.Header.Add("Authorization", "Bearer "+os.Getenv("gitlab_token"))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("getExternalUsernames res error: ", err.Error())
			break
		}
		if res.StatusCode != 200 {
			fmt.Println("getExternalUsernames status code error: ", res.StatusCode)
			break
		}
		resBody, err := ioutil.ReadAll(res.Body)
		ud := []UserDetail{}
		err = json.Unmarshal(resBody, &ud)
		if err != nil {
			fmt.Println("getExternalUsernames unmashal error: ", err.Error())
			break
		}
		if reflect.DeepEqual(ud, []UserDetail{}) {
			break
		}
		for _, v := range ud {
			externalUsernames = append(externalUsernames, v.Username)
		}
		page++
	}
	return externalUsernames
}

// It returns developers, maintainers and owners usernames
func GetProjectMembers(projectId string) ([]string, []string) {
	others := []string{}
	owners := []string{}
	excludedUsernames := GetExcludedUsernames()
	externalUsernames := GetExternalUsernames()
	page := 1
	for {
		req, err := http.NewRequest("GET", os.Getenv("gitlab_api_url")+"projects/"+
			projectId+"/members/all?page="+
			strconv.Itoa(page)+"&per_page=100", nil)
		if err != nil {
			fmt.Println("getProjectMembers req error: ", err.Error())
			break
		}
		req.Header.Add("Authorization", "Bearer "+os.Getenv("gitlab_token"))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("getProjectMembers res error: ", err.Error())
			break
		}
		if res.StatusCode != 200 {
			fmt.Println("getProjectMembers status code error: ", res.StatusCode)
			break
		}
		resBody, err := ioutil.ReadAll(res.Body)
		ud := []UserDetail{}
		err = json.Unmarshal(resBody, &ud)
		if err != nil {
			fmt.Println("getProjectMembers unmashal error: ", err.Error())
			break
		}
		if reflect.DeepEqual(ud, []UserDetail{}) {
			break
		}
		for _, v := range ud {
			if ContainsInSclice(excludedUsernames, v.Username) || ContainsInSclice(externalUsernames, v.Username) {
				continue
			}
			if !strings.Contains(v.Username, "_bot") && v.State == "active" {
				if v.AccessLevel == 30 || v.AccessLevel == 40 {
					others = append(others, v.Username)
				} else if v.AccessLevel == 50 {
					owners = append(owners, v.Username)
				}
			}
		}
		page++
	}

	return others, owners
}

// It returns owners and assignee username
func GetWatchersAndAssignee(projectId string, username string) ([]string, string) {
	others, owners := GetProjectMembers(projectId)
	if ContainsInSclice(others, username) {
		others = RemoveInSlice(others, username)
	} else if ContainsInSclice(owners, username) {
		owners = RemoveInSlice(owners, username)
	}
	ShuffleSlice(others)
	ShuffleSlice(owners)
	if len(owners) > 5 {
		owners = owners[:5]
	}
	assignee := ""
	if len(others) > 0 {
		assignee = others[0]
		others = others[1:]
	}
	if len(others) > 0 && len(owners) < 5 {
		owners = append(owners, others...)
		if len(owners) > 5 {
			owners = owners[:5]
		}
	}
	return owners, assignee
}

// It finds username by an email.
func GetUserByEmail(email string) (string, string) {
	username := ""
	name := ""
	req, err := http.NewRequest("GET", os.Getenv("gitlab_api_url")+"users?search="+email, nil)
	if err != nil {
		fmt.Println("GetUserByEmail: req error: ", err)
		return username, name
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("gitlab_token"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("GetUserByEmail: res error: ", err.Error())
		return username, name
	}

	if res.StatusCode == 200 {
		resBody, _ := ioutil.ReadAll(res.Body)
		ud := []UserDetail{}
		err := json.Unmarshal(resBody, &ud)
		if err != nil {
			fmt.Println("GetUserByEmail: unmashal error: ", err.Error())
			return username, name
		}
		if reflect.DeepEqual(ud, []UserDetail{}) {
			return username, name
		}
		username = ud[0].Username
		name = ud[0].Name
	}
	return username, name
}
