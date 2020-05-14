package eureka

import (
	"strconv"
)

type status struct {
	Name   string `json:"name"`
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
}

type health struct {
	Status  string  `json:"status"`
	Details Details `json:"details"`
}

type Details struct {
}

func (client *Client) ActuatorStatus() interface{} {
	appStatus := status{}
	appStatus.Name = client.config.InstanceConfig.AppName
	appStatus.Server.Port = strconv.Itoa(client.GetPort())
	return appStatus
}

func (client *Client) ActuatorHealth() interface{} {
	appHealth := health{}
	appHealth.Status = "UP"
	appHealth.Details = Details{}
	return appHealth
}
