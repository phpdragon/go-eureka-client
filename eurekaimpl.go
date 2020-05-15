package eureka

import (
	"fmt"
	"github.com/phpdragon/go-eureka-client/core"
	netUtil "github.com/phpdragon/go-eureka-client/netutil"
	"time"
)

// Api for sending rest httpClient to eureka server
func (client *Client) Api() (*core.EurekaServerApi, error) {
	api, err := client.pickEurekaServerApi()
	if err != nil {
		return nil, err
	}
	return api, nil
}

func (client *Client) refreshRegistry() {
	if !client.config.ClientConfig.FetchRegistry {
		return
	}

	for {
		_ = client.fetchRegistry()
		time.Sleep(time.Second * time.Duration(client.config.ClientConfig.GetRegistryFetchIntervalSeconds()))
	}
}

//刷新服务列表
func (client *Client) fetchRegistry() error {
	client.logger.Info("Fetch registry info")

	api, err := client.Api()
	if err != nil {
		client.logger.Error(fmt.Sprintf("Failed to QueryAllInstances, err=%s", err.Error()))
		return err
	}

	apps, err := api.QueryAllInstances()
	if err != nil {
		client.logger.Error(fmt.Sprintf("Failed to QueryAllInstances, err=%s", err.Error()))
		return err
	}

	registryApps := make(map[string]*core.Application)
	activeInstances := make(map[string]map[int]*core.Instance)
	activeServiceUrls := make(map[string]map[int]map[int]string)

	for _, app := range apps.Applications {
		instances, urls := getActiveInstancesAndIpPorts(client.config.ClientConfig.FilterOnlyUpInstances, app.Instances)
		registryApps[app.Name] = &app
		activeInstances[app.Name] = instances
		activeServiceUrls[app.Name] = urls
	}

	client.mutex.Lock()
	defer client.mutex.Unlock()

	client.registryAppMap = registryApps
	client.activeInstanceMap = activeInstances
	client.activeServiceIpPortMap = activeServiceUrls

	return nil
}

// register instance (default current status is STARTING)
// and update instance status to UP
func (client *Client) registerWithEureka() {
	if !client.config.ClientConfig.RegisterWithEureka {
		client.logger.Warn("This instance don't register to eureka!")
		return
	}

	for {
		if client.instance == nil {
			client.logger.Error("Config instance can't be nil")
			return
		}

		api, err := client.Api()
		if err != nil {
			time.Sleep(time.Second * defaultSleepIntervals)
			continue
		}

		err = api.RegisterInstance(client.instance.App, client.instance)
		if err != nil {
			client.logger.Error(fmt.Sprintf("clientConfig register failed, err=%s", err.Error()))
			time.Sleep(time.Second * defaultSleepIntervals)
			continue
		}
		client.logger.Info(fmt.Sprintf("Successfully register service to eureka with status[%s] !", client.instance.Status))

		break
	}

	go func() {
		for {
			enabledOnInit := client.config.InstanceConfig.InstanceEnabledOnInit
			//如果向eureka注册后立即启用实例以获取流量，或者服务已经启动，则向eureka更新为在线状态
			if enabledOnInit || (!enabledOnInit && client.serverIsStarted()) {
				updated, err := client.updateInstanceStatus()
				if nil != err {
					client.logger.Error(err.Error())
				}
				if updated {
					break
				}
			}
			time.Sleep(time.Second * defaultSleepIntervals)
		}
	}()

	//发送心跳
	go client.heartbeat()
}

//判断http服务是否已经启动
func (client *Client) serverIsStarted() bool {
	port := client.instance.Port.Port
	if "true" == client.instance.SecurePort.Enabled {
		port = client.instance.SecurePort.Port
	}

	used := netUtil.PortInUse(client.instance.IpAddr, port)
	client.logger.Debug(fmt.Sprintf("Check that the web server is started, result:%t", used))

	return used
}

func (client *Client) updateInstanceStatus() (bool, error) {
	client.logger.Info("Update the instance status to UP ...")

	if client.instance == nil {
		client.logger.Error("Config instance can't be nil")
		return false, nil
	}

	api, err := client.Api()
	if err != nil {
		return false, nil
	}

	//如果成功注册到eureka并将状态更新到UP
	// if success to register to eureka and update status to UP
	// then break loop
	err = api.UpdateInstanceStatus(client.instance.App, client.instance.InstanceId, core.STATUS_UP)
	if err != nil {
		client.logger.Error(fmt.Sprintf("clientConfig UP failed, err=%s", err.Error()))
		return false, nil
	}

	//本地状态更新为up
	client.instance.Status = core.STATUS_UP

	client.logger.Info("The server status[UP] was updated successfully !")

	return true, nil
}

//TODO:
// rand to pick service url and new EurekaServerApi instance
func (client *Client) pickEurekaServerApi() (*core.EurekaServerApi, error) {
	return core.NewEurekaServerApi(client.config.ServiceURL.DefaultZone), nil
}

// 发送心跳
// eureka client heartbeat
func (client *Client) heartbeat() {
	for {
		api, err := client.Api()
		if err != nil {
			time.Sleep(time.Second * defaultSleepIntervals)
			continue
		}

		err = api.SendHeartbeat(client.instance.App, client.instance.InstanceId)
		if err != nil {
			client.logger.Error(fmt.Sprintf("Failed to send heartbeat, err=%s", err.Error()))
			time.Sleep(time.Second * defaultSleepIntervals)
			continue
		}

		client.logger.Debug(fmt.Sprintf("Heartbeat app=%s, instanceId=%s", client.instance.App, client.instance.InstanceId))
		time.Sleep(time.Duration(client.config.InstanceConfig.LeaseInfo.RenewalIntervalInSecs) * time.Second)
	}
}

func  (client *Client) monitor() {

}