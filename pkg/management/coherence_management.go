/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// The management package contains types and functions for working with Coherence management over ReST.
package management

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// The URL pattern for Coherence management cluster query.
	clusterFormat = "http://%s:%d/management/coherence/cluster"
	// The URL pattern for Coherence management members query.
	membersFormat = "http://%s:%d/management/coherence/cluster/members"
	// The URL pattern for Coherence management services query.
	servicesFormat = "http://%s:%d/management/coherence/cluster/services"
	// The URL pattern for Coherence management partition assignment query.
	partitionFormat = "http://%s:%d/management/coherence/cluster/services/%s/partition"
)

// A struct to use to hold the results of a generic Coherence management ReST query.
type RestData struct {
	Links []map[string]string
	Items []map[string]interface{}
}

// A struct to use to hold the results of a generic Coherence management ReST cluster query.
type ClusterData struct {
	Links         []map[string]string `json:"Links"`
	RefreshTime   string              `json:"refreshTime"`
	LicenseMode   string              `json:"licenseMode"`
	ClusterSize   int                 `json:"clusterSize"`
	LocalMemberID int                 `json:"localMemberId"`
	Version       string              `json:"version"`
	Running       bool                `json:"running"`
	ClusterName   string              `json:"clusterName"`
}

// A struct to use to hold the results of a Coherence management ReST services query
// http://localhost:30000/management/coherence/cluster/services
type ServicesData struct {
	Links []map[string]string `json:"Links"`
	Items []ServiceData
}

// A struct to use to hold the results of a Coherence management ReST service query
// http://localhost:30000/management/coherence/cluster/services/%s
type ServiceData struct {
	Links []map[string]string `json:"Links"`
	Name  string              `json:"name"`
	Type  string              `json:"type"`
}

// A struct to use to hold the results of a Coherence management ReST partition assignment query
// http://localhost:30000/management/coherence/cluster/services/%s/partition
// This structure only contains a sub-set of the fields available in the response json. If other
// fields are required they should be added to this struct.
type PartitionData struct {
	Links                      []map[string]string `json:"Links"`
	HAStatus                   string              `json:"HAStatus"`
	HAStatusCode               int                 `json:"HAStatusCode"`
	RemainingDistributionCount int                 `json:"remainingDistributionCount"`
	BackupCount                int                 `json:"backupCount"`
	ServiceNodeCount           int                 `json:"serviceNodeCount"`
}

// A struct to use to hold the results of a Coherence management ReST members query
// http://localhost:30000/management/coherence/cluster/members
type MembersData struct {
	Links []map[string]string `json:"Links"`
	Items []MemberData
}

// A struct to use to hold the results of a Coherence management ReST member query.
// http://localhost:30000/management/coherence/cluster/members/<member-id>
// This structure only contains a sub-set of the fields available in the response json. If other
// fields are required they should be added to this struct.
type MemberData struct {
	Links        []map[string]string `json:"Links"`
	SiteName     string              `json:"siteName"`
	RackName     string              `json:"rackName"`
	MachineName  string              `json:"machineName"`
	MachineID    int                 `json:"machineId"`
	MemberName   string              `json:"memberName"`
	RoleName     string              `json:"roleName"`
	ID           int                 `json:"id"`
	NodeID       string              `json:"nodeId"`
	LoggingLevel int                 `json:"loggingLevel"`
}

// Perform a Management over ReST cluster query http://localhost:30000/management/coherence/cluster
// and return the results, the http response status and any error.
func GetCluster(cl *http.Client, host string, port int32) (*ClusterData, int, error) {
	url := fmt.Sprintf(clusterFormat, host, port)
	data := &ClusterData{}
	status, err := query(cl, url, data)
	return data, status, err
}

// Perform a Management over ReST members query http://localhost:30000/management/coherence/cluster/members
// and return the results, the http response status and any error.
func GetMembers(cl *http.Client, host string, port int32) (*MembersData, int, error) {
	url := fmt.Sprintf(membersFormat, host, port)
	data := &MembersData{}
	status, err := query(cl, url, data)
	return data, status, err
}

// Perform a Management over ReST members query http://localhost:30000/management/coherence/cluster/services
// and return the results, the http response status and any error.
func GetServices(cl *http.Client, host string, port int32) (*ServicesData, int, error) {
	url := fmt.Sprintf(servicesFormat, host, port)
	data := &ServicesData{}
	status, err := query(cl, url, data)
	return data, status, err
}

// Perform a Management over ReST members query http://localhost:30000/management/coherence/cluster/services/%s/partition
// and return the results, the http response status and any error.
func GetPartitionAssignment(cl *http.Client, host string, port int32, service string) (*PartitionData, int, error) {
	url := fmt.Sprintf(partitionFormat, host, port, service)
	data := &PartitionData{}
	status, err := query(cl, url, data)
	return data, status, err
}

// Perform a Management over ReST query and parse the json response if the response code is 200
// returning the response code any any error.
func query(cl *http.Client, url string, v interface{}) (int, error) {
	var response *http.Response
	var err error

	// re-try a max of 5 times
	for i := 0; i < 5; i++ {
		response, err = cl.Get(url)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		var status = http.StatusInternalServerError
		if response != nil {
			status = response.StatusCode
		}
		return status, err
	}

	if response.StatusCode == http.StatusOK {
		data, _ := ioutil.ReadAll(response.Body)

		err = json.Unmarshal(data, v)
	}

	return response.StatusCode, err
}
