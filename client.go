package eureka

import (
	"fmt"
	"github.com/phpdragon/go-eureka-client/config"
	"github.com/phpdragon/go-eureka-client/core"
	"github.com/phpdragon/go-eureka-client/logger"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	defaultSleepIntervals = 3
	//
	httpPrefix  = "http://"
	httpsPrefix = "https://"
	//
	httpKey  = 0
	httpsKey = 1
)

type Client struct {
	Running bool

	//
	apiClient *core.EurekaServerApi

	//自增器
	autoIncr *atomic.Int64

	// for monitor system signal
	signalChan chan os.Signal

	//日志对象
	logger *logger.Logger

	mutex sync.RWMutex

	config *config.Config

	// current client (instance) config
	instance *core.Instance

	// applications registry
	// key: appId
	// value: Application
	registryAppMap map[string]*core.Application

	// instances registry
	// key: appId
	// value:
	//		key:  int(0...n)
	//		value: InstanceConfig
	activeInstanceMap map[string]map[int]*core.Instance

	// instance real url map
	// key: appId
	// value:
	//		key:  int(http:0, https:1)
	//		value:
	//			key:  int(0...n)
	//			value: real url
	activeServiceIpPortMap map[string]map[int]map[int]string
}

func NewClient(configPath string) *Client {
	return NewClientWithLog(configPath, nil)
}

func NewClientWithLog(configPath string, zapLog *zap.Logger) *Client {
	logger := logger.NewLogAgent(zapLog)

	eurekaConfig, err := config.LoadConfig(configPath, false)
	if err != nil {
		logger.Error(fmt.Sprintf("LoadConfig %s failed, err=%s", configPath, err.Error()))
		os.Exit(1)
	}

	//实例化
	instanceInfo, err := config.NewInstance(eurekaConfig)
	if err != nil {
		logger.Error(fmt.Sprintf("NewInstance %s failed, err=%s", eurekaConfig.InstanceConfig.AppName, err.Error()))
		os.Exit(1)
	}

	client := &Client{
		//自增器
		autoIncr:   atomic.NewInt64(0),
		logger:     logger,
		signalChan: make(chan os.Signal),
		//
		config:   eurekaConfig,
		instance: instanceInfo,
	}

	api, err := client.Api()
	if err != nil {
		client.logger.Error(fmt.Sprintf("Failed to get EurekaServerApi instance, err=%s", err.Error()))
		os.Exit(1)
	}
	client.apiClient = api

	return client
}

func (client *Client) Run() {
	client.mutex.Lock()
	client.Running = true
	client.mutex.Unlock()

	// handle exit signal to de-register instance
	go client.handleSignal()

	// (if FetchRegistry is true), fetch registry apps periodically
	// and update to t.registryAppMap
	go client.refreshRegistry()

	client.registerWithEureka()
}

func (client *Client) Shutdown() {
	//client在shutdown情况下，是否显示从注册中心注销
	if !client.Running || !client.config.ClientConfig.ShouldUnregisterOnShutdown {
		return
	}

	client.logger.Info(fmt.Sprintf("Receive exit signal, client instance going to de-register, instanceId=%s.", client.instance.InstanceId))

	// de-register instance
	err := client.apiClient.DeRegisterInstance(client.instance.App, client.instance.InstanceId)
	if err != nil {
		client.logger.Error(fmt.Sprintf("Failed to de-register %s, err=%s", client.instance.InstanceId, err.Error()))
		return
	}

	client.mutex.Lock()
	client.Running = false
	client.mutex.Unlock()

	client.logger.Info(fmt.Sprintf("de-register %s success.", client.instance.InstanceId))
}

// for graceful kill. Here handle SIGTERM signal to do sth
// e.g: kill -TERM $pid
//      or "ctrl + c" to exit
func (client *Client) handleSignal() {
	if client.signalChan == nil {
		client.signalChan = make(chan os.Signal)
	}

	signal.Notify(client.signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		switch <-client.signalChan {
		case syscall.SIGINT:
			client.logger.Info(fmt.Sprintf("syscall.SIGINT, instanceId=%s.", client.instance.InstanceId))
			fallthrough
		case syscall.SIGKILL:
			client.logger.Info(fmt.Sprintf("syscall.SIGKILL, instanceId=%s.", client.instance.InstanceId))
			fallthrough
		case syscall.SIGHUP:
			client.logger.Info(fmt.Sprintf("syscall.SIGHUP, instanceId=%s.", client.instance.InstanceId))
			fallthrough
		case syscall.SIGQUIT:
			client.logger.Info(fmt.Sprintf("syscall.SIGQUIT, instanceId=%s.", client.instance.InstanceId))
			fallthrough
		case syscall.SIGTERM:
			client.Shutdown()
			os.Exit(0)
		}
	}
}
