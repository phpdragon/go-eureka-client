package eureka

import (
	"fmt"
	"github.com/phpdragon/go-eurake-client/core"
	"github.com/phpdragon/go-eurake-client/misc"
	"net"
	"strings"
	"time"
)

//获取下一个容器
func (client *Client) GetNextServerFromEureka(appId string) (*core.Instance, error) {
	instanceMap, err := client.getActiveInstancesByAppId(appId)
	if nil != err {
		return &core.Instance{}, err
	}

	if nil == instanceMap || 0 == len(instanceMap) {
		client.logger.Error(fmt.Sprintf("This %s instances not exist!", appId))
		return &core.Instance{}, fmt.Errorf("This %s instances not exist!", appId)
	}

	index := client.getRandIndex(len(instanceMap))
	return instanceMap[index], nil
}

func (client *Client) GetRealHttpUrl(httpUrl string) (string, error) {
	httpUrlTmp := strings.Replace(httpUrl, httpPrefix, "", -1)
	httpUrlTmp = strings.Replace(httpUrlTmp, httpsPrefix, "", -1)
	urls := strings.Split(httpUrlTmp, "/")
	appName := urls[0]

	//是否https
	mapKey := httpKey
	if strings.HasPrefix(httpUrl, httpsPrefix) {
		mapKey = httpsKey
	}

	ipPortMap, err := client.getActiveServiceIpPortByAppId(appName)
	if nil != err || 0 == len(ipPortMap) {
		//TODO：文案
		return "", fmt.Errorf("This %s instances not exist!", appName)
	}

	//取http还是https的ip:port
	realIpPorts := ipPortMap[mapKey]
	if nil == realIpPorts || 0 == len(realIpPorts) {
		//TODO：文案
		return "", fmt.Errorf("This %s instances not exist!", appName)
	}

	//随机取一个目标ip:port
	index := client.getRandIndex(len(realIpPorts))
	realIpPort := realIpPorts[index]

	return strings.Replace(httpUrl, appName, realIpPort, -1), nil
}

func (client *Client) GetApplications() map[string]*core.Application {
	return client.registryAppMap
}

func (client *Client) GetInstances() map[string]map[int]*core.Instance {
	return client.activeInstanceMap
}

func (client *Client) ActuatorStatus() interface{} {
	return misc.ActuatorStatus(client.GetPort(), client.GetAppName())
}

func (client *Client) ActuatorHealth() interface{} {
	return misc.ActuatorHealth()
}

func (client *Client) GetAppName() string {
	return client.config.InstanceConfig.AppName
}

func (client *Client) GetPort() int {
	port := client.instance.Port.Port
	if "true" == client.instance.SecurePort.Enabled {
		port = client.instance.SecurePort.Port
	}
	return port
}

func (client *Client) PortInUse(host string, ports []string) bool {
	for _, port := range ports {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err != nil {
			return false
		}
		if conn != nil {
			defer conn.Close()
			return true
		}
	}

	return false
}
