package sls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type ETLJobV2 struct {
	Configuration    ETLConfiguration `json:"configuration"`
	Description      string           `json:"description"`
	DisplayName      string           `json:"displayName"`
	Name             string           `json:"name"`
	Schedule         ETLSchedule      `json:"schedule"`
	Type             string           `json:"type"`
	Status           string           `json:"status"`
	CreateTime       int32            `json:"createTime"`
	LastModifiedTime int32            `json:"lastModifiedTime"`
}

type ETLConfiguration struct {
	AccessKeyId     string            `json:"accessKeyId"`
	AccessKeySecret string            `json:"accessKeySecret"`
	FromTime        int64             `json:"fromTime"`
	Logstore        string            `json:"logstore"`
	Parameters      map[string]string `json:"parameters"`
	RoleArn         string            `json:"roleArn"`
	Script          string            `json:"script"`
	ToTime          int32             `json:"toTime"`
	Version         int8              `json:"version"`
	ETLSinks        []ETLSink         `json:"sinks"`
}

type ETLSchedule struct {
	Type string `json:"type"`
}

type ETLSink struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	Endpoint        string `json:"endpoint"`
	Logstore        string `json:"logstore"`
	Name            string `json:"name"`
	Project         string `json:"project"`
	RoleArn         string `json:"roleArn"`
}

type ListETLResponse struct {
	Total   int         `json:"total"`
	Count   int         `json:"count"`
	Results []*ETLJobV2 `json:"results"`
}


func NewLogETLJobV2(endpoint, accessKeyId, accessKeySecret, logstore, name, project string) ETLJobV2 {
	sink := ETLSink{
		AccessKeyId:accessKeyId,
		AccessKeySecret:accessKeySecret,
		Endpoint:endpoint,
		Logstore:logstore,
		Name:name,
		Project:project,
	}
	config := ETLConfiguration {
		AccessKeyId:accessKeyId,
		AccessKeySecret:accessKeySecret,
		FromTime: time.Now().Unix(),
		Script: "e_set('new','aliyun')",
		Version:2,
		Logstore:logstore,
		ETLSinks:[]ETLSink{sink},
		Parameters: map[string]string{},

	}
	schedule := ETLSchedule{
		Type:"Resident",
	}
	etljob := ETLJobV2 {
		Configuration:config,
		DisplayName:"displayname",
		Description:"go sdk case",
		Name:name,
		Schedule:schedule,
		Type:"ETL",

	}
	return etljob
}



func (c *Client) CreateETL(project string, etljob ETLJobV2) error {
	body, err := json.Marshal(etljob)
	if err != nil {
		return NewClientError(err)
	}
	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}
	uri := "/jobs"

	r, err := c.request(project, "POST", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) GetETL(project string, etlName string) (ETLJob *ETLJobV2, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	uri := "/jobs/" + etlName
	r, err := c.request(project, "GET", uri, h, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	buf, _ := ioutil.ReadAll(r.Body)
	etlJob := &ETLJobV2{}
	if err = json.Unmarshal(buf, etlJob); err != nil {
		err = NewClientError(err)
	}
	return etlJob, nil
}

func (c *Client) UpdateETL(project string, etljob ETLJobV2) error {
	body, err := json.Marshal(etljob)
	if err != nil {
		return NewClientError(err)
	}
	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}
	uri := "/jobs/" + etljob.Name
	r, err := c.request(project, "PUT", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) DeleteETL(project string, etlName string) error {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	uri := "/jobs/" + etlName
	r, err := c.request(project, "DELETE", uri, h, nil)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) ListETL(project string, offset int, size int) (*ListETLResponse, error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}

	uri := fmt.Sprintf("/jobs?offset=%d&size=%d", offset, size)
	r, err := c.request(project, "GET", uri, h, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	buf, _ := ioutil.ReadAll(r.Body)

	listETLResponse := &ListETLResponse{}
	if err = json.Unmarshal(buf, listETLResponse); err != nil {
		err = NewClientError(err)
	}
	return listETLResponse, err
}

func (c *Client) StartETL(project, name string) error {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}

	uri := fmt.Sprintf("/jobs/%s?action=START", name)
	r, err := c.request(project, "PUT", uri, h, nil)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) StopETL(project, name string) error {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}

	uri := fmt.Sprintf("/jobs/%s?action=STOP", name)
	fmt.Println(uri)
	r, err := c.request(project, "PUT", uri, h, nil)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}
