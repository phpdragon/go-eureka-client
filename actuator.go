package eureka

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
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
	Details details `json:"details"`
}

type details struct {
}

type href struct {
	Href      string `json:"href"`
	Templated bool   `json:"templated"`
}

type routeInfo struct {
	path     string
	function handlerFunc
}

type handlerFunc func(client *Client) interface{}

var routePath = []routeInfo{
	//处理eureka的心跳等
	{path: "^/actuator$", function: actuatorLinks},
	{path: "^/actuator/info$", function: actuatorInfo},
	{path: "^/actuator/health$", function: actuatorHealth},
	{path: "^/actuator\\w*", function: actuatorAny},
}

func (client *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range routePath {
		ok, _ := regexp.Match(route.path, []byte(r.URL.Path))
		if ok {
			client.writeJson(w, r, route.function(client), true)
			return
		}
	}
	_, _ = w.Write([]byte("404 not found"))
}

func (client *Client) writeJson(rw http.ResponseWriter, req *http.Request, response interface{}, isJson bool) {
	origin := req.Header.Get("origin")
	rw.Header().Set("cache-control", "No-Cache")
	rw.Header().Set("content-type", "application/json; charset=utf-8")
	rw.Header().Set("pragma", "No-Cache")
	rw.Header().Set("expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	if 0 < len(origin) {
		rw.Header().Set("access-control-allow-origin", origin)
		rw.Header().Set("access-control-allow-credentials", "true")
	}

	var err error
	var dataBody []byte
	if isJson {
		dataBody, err = client.toStringByte(response)
		if err != nil {
			client.logger.Error(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		dataBody = []byte(response.(string))
	}

	_, err = rw.Write(dataBody)
	if err != nil {
		client.logger.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (client *Client) toStringByte(v interface{}) ([]byte, error) {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		client.logger.Error("解析失败:", err.Error())
		return nil, nil
	}
	return jsonStr, nil
}

func actuatorLinks(client *Client) interface{} {
	links := make(map[string]href, 10)

	url := fmt.Sprintf("http://%s:%s", client.instance.IpAddr, client.config.InstanceConfig.NonSecurePort)
	links["info"] = href{
		Href:      url + "/actuator/info",
		Templated: false,
	}
	links["health"] = href{
		Href:      url + "/actuator/health",
		Templated: false,
	}
	return links
}

func actuatorInfo(client *Client) interface{} {
	appStatus := status{}
	appStatus.Name = client.config.InstanceConfig.AppName
	appStatus.Server.Port = strconv.Itoa(client.GetPort())
	return appStatus
}

func actuatorHealth(*Client) interface{} {
	appHealth := health{}
	appHealth.Status = "UP"
	appHealth.Details = details{}
	return appHealth
}

func actuatorAny(*Client) interface{} {
	return new(interface{})
}
