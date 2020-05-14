package eureka

import (
	"fmt"
	"github.com/phpdragon/go-eurake-client/core"
	"strings"
	"unsafe"
)

func (client *Client) getActiveServiceIpPortByAppId(appId string) (map[int]map[int]string, error) {
	id := strings.ToUpper(appId)
	cache := client.activeServiceIpPortMap[id]
	if nil != cache {
		return client.activeServiceIpPortMap[id], nil
	}

	err := client.doRefreshByAppId(appId)
	if nil != err {
		return nil, err
	}

	return client.activeServiceIpPortMap[id], nil
}

func (client *Client) getActiveInstancesByAppId(appId string) (map[int]*core.Instance, error) {
	id := strings.ToUpper(appId)
	cache := client.activeInstanceMap[id]
	if nil != cache {
		return client.activeInstanceMap[id], nil
	}

	err := client.doRefreshByAppId(appId)
	if nil != err {
		return nil, err
	}

	return client.activeInstanceMap[id], nil
}

func (client *Client) doRefreshByAppId(appId string) error {
	api, err := client.Api()
	if err != nil {
		return err
	}

	application, errr := api.QueryAllInstanceByAppId(appId)
	if errr != nil {
		return errr
	}

	instances, urls := getActiveInstancesAndIpPorts(client.config.ClientConfig.FilterOnlyUpInstances, application.Instances)

	client.mutex.Lock()
	defer client.mutex.Unlock()

	client.registryAppMap[appId] = application
	client.activeInstanceMap[appId] = instances
	client.activeServiceIpPortMap[appId] = urls

	return nil
}

//获取有效的实例和链接
func getActiveInstancesAndIpPorts(filterOnlyUpInstances bool, instances []core.Instance) (map[int]*core.Instance, map[int]map[int]string) {
	instancesX := make(map[int]*core.Instance)
	//
	urls := make(map[int]map[int]string)
	httpUrls := make(map[int]string)
	httpsUrls := make(map[int]string)
	instanceTotal := len(instances)
	for i := 0; i < instanceTotal; i++ {
		instance := instances[i]

		if filterOnlyUpInstances && core.STATUS_UP != instance.Status {
			continue
		}

		instancesX[i] = &instance

		if "true" == instance.Port.Enabled {
			httpUrls[i] = fmt.Sprintf("%s:%d", instance.IpAddr, instance.Port.Port)
		}
		if "false" == instance.SecurePort.Enabled {
			httpsUrls[i] = fmt.Sprintf("%s:%d", instance.IpAddr, instance.SecurePort.Port)
		}
	}

	urls[httpKey] = httpUrls
	urls[httpsKey] = httpsUrls
	return instancesX, urls
}

func (client *Client) getRandIndex(total int) int {
	var index64 = client.autoInc.Inc() % int64(total)
	return *(*int)(unsafe.Pointer(&index64))
}
