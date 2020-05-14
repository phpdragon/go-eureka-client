package eureka

import (
	"bytes"
	"errors"
	"fmt"
	yaml "github.com/go-yaml/yaml"
	core "github.com/phpdragon/go-eurake-client/core"
	netUtil "github.com/phpdragon/go-eurake-client/netutil"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	template "text/template"
)

type templateData struct {
	Env map[string]string
}

type inValid struct {
	Name string
	Type string
}

type (
	Root struct {
		Server struct {
			Port int `yaml:"port"`
		} `yaml:"server"`
		Eureka Config `yaml:"eureka"`
	}

	//see: https://www.cnblogs.com/liukaifeng/p/10052594.html
	Config struct {
		ServiceURL struct {
			DefaultZone string `yaml:"defaultZone"`
		} `yaml:"serviceUrl"`
		ClientConfig   ClientConfig `yaml:"client"`
		InstanceConfig struct {
			InstanceId            string `yaml:"instanceId"`
			AppName               string `yaml:"appName"`
			NonSecurePort         int    `yaml:"nonSecurePort"`
			NonSecurePortEnabled  bool   `yaml:"nonSecurePortEnabled"`
			SecurePort            int    `yaml:"securePort"`
			SecurePortEnabled     bool   `yaml:"securePortEnabled"`
			SecureVirtualHostName string `yaml:"secureVirtualHostName"`
			VirtualHostName       string `yaml:"virtualHostName"`
			//
			HomePageUrlPath    string `yaml:"homePageUrlPath"`
			StatusPageUrlPath  string `yaml:"statusPageUrlPath"`
			HealthCheckUrlPath string `yaml:"healthCheckUrlPath"`
			//
			PreferIpAddress       bool `yaml:"preferIpAddress"`
			InstanceEnabledOnInit bool `yaml:"instanceEnabledOnInit"`
			//
			CountryId int                    `yaml:"countryId"`
			Metadata  map[string]interface{} `yaml:"metadata"`
			LeaseInfo struct {
				RenewalIntervalInSecs int `yaml:"renewalIntervalInSecs"`
				DurationInSecs        int `yaml:"durationInSecs"`
			} `yaml:"leaseInfo"`
		} `yaml:"instance"`
	}

	ClientConfig struct {
		//指示从eureka服务器获取注册表信息的频率,默认30s
		RegistryFetchIntervalSeconds int `yaml:"registryFetchIntervalSeconds"`
		//客户端是否获取eureka服务器注册表上的注册信息,不调用其他微服务可以为false，默认为false
		FetchRegistry bool `yaml:"fetchRegistry"`
		//是否过滤掉非up实例，默认为false
		FilterOnlyUpInstances bool `yaml:"filterOnlyUpInstances"`
		//指示此实例是否应将其信息注册到eureka服务器以供其他服务发现，默认为false
		RegisterWithEureka bool `yaml:"registerWithEureka"`
		//client在shutdown情况下，是否显示从注册中心注销
		ShouldUnregisterOnShutdown bool `yaml:"shouldUnregisterOnShutdown"`
	}
)

func LoadConfig(configPath string, valid bool) (*Config, error) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return &Config{}, err
	}

	file, err = substitute(file)
	if err != nil {
		return &Config{}, err
	}

	config := &Root{}
	if err = yaml.Unmarshal(file, config); err != nil {
		return &Config{}, err
	}

	if valid {
		if err = validate(config); err != nil {
			return &Config{}, err
		}
	}

	return &config.Eureka, nil
}

func NewInstance(config *Config) (*core.Instance, error) {
	if isEmpty(config.InstanceConfig.AppName) {
		return nil, fmt.Errorf("eureka.instance.appName is empty！")
	}

	//是否优先使用服务实例的IP地址，相较于
	var localIp = ""
	if config.InstanceConfig.PreferIpAddress {
		localIp = netUtil.GetLocalIp()
	} else {
		hostname, err := os.Hostname()
		if nil != err {
			return &core.Instance{}, err
		}

		localIp = hostname
	}

	instance := &core.Instance{
		InstanceId:       config.InstanceConfig.InstanceId,
		HostName:         localIp,
		IpAddr:           localIp,
		App:              config.InstanceConfig.AppName,
		VipAddress:       config.InstanceConfig.VirtualHostName,
		SecureVipAddress: config.InstanceConfig.SecureVirtualHostName,
		Status:           core.STATUS_STARTING,
		Port: &core.Port{
			Port:    config.InstanceConfig.NonSecurePort,
			Enabled: strconv.FormatBool(config.InstanceConfig.NonSecurePortEnabled),
		},
		SecurePort: &core.Port{
			Port:    config.InstanceConfig.SecurePort,
			Enabled: strconv.FormatBool(config.InstanceConfig.SecurePortEnabled),
		},
		HomePageUrl:    config.InstanceConfig.HomePageUrlPath,
		StatusPageUrl:  config.InstanceConfig.StatusPageUrlPath,
		HealthCheckUrl: config.InstanceConfig.HealthCheckUrlPath,
		// 数据中心
		DataCenterInfo: &core.DataCenterInfo{
			Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
			Name:  core.DC_NAME_TYPE_MY_OWN,
			//Metadata: &core.DataCenterMetadata{
			//	AmiLaunchIndex:   config.InstanceConfig.DataCenterInfo.Metadata.AmiLaunchIndex,
			//	LocalHostname:    config.InstanceConfig.DataCenterInfo.Metadata.LocalHostname,
			//	AvailabilityZone: config.InstanceConfig.DataCenterInfo.Metadata.AvailabilityZone,
			//	InstanceId:       config.InstanceConfig.DataCenterInfo.Metadata.InstanceId,
			//	PublicIpv4:       config.InstanceConfig.DataCenterInfo.Metadata.PublicIpv4,
			//	PublicHostname:   config.InstanceConfig.DataCenterInfo.Metadata.PublicHostname,
			//	AmiManifestPath:  config.InstanceConfig.DataCenterInfo.Metadata.AmiManifestPath,
			//	LocalIpv4:        config.InstanceConfig.DataCenterInfo.Metadata.LocalIpv4,
			//	Hostname:         config.InstanceConfig.DataCenterInfo.Metadata.Hostname,
			//	AmiID:            config.InstanceConfig.DataCenterInfo.Metadata.AmiID,
			//	InstanceType:     config.InstanceConfig.DataCenterInfo.Metadata.InstanceType,
			//},
		},
		Metadata: config.InstanceConfig.Metadata,
		LeaseInfo: &core.LeaseInfo{
			RenewalIntervalInSecs: config.InstanceConfig.LeaseInfo.RenewalIntervalInSecs,
			DurationInSecs:        config.InstanceConfig.LeaseInfo.DurationInSecs,
		},
		CountryID: config.InstanceConfig.CountryId,
	}

	dataCenterName := instance.DataCenterInfo.Name
	if core.DC_NAME_TYPE_MY_OWN != dataCenterName && core.DC_NAME_TYPE_AMAZON != dataCenterName {
		instance.DataCenterInfo.Name = core.DC_NAME_TYPE_MY_OWN
	}

	port := config.InstanceConfig.NonSecurePort
	if config.InstanceConfig.SecurePortEnabled {
		port = config.InstanceConfig.SecurePort
	}

	if isEmpty(config.InstanceConfig.InstanceId) {
		instance.InstanceId = fmt.Sprintf("%s:%d", localIp, port)
	}
	if isEmpty(instance.HostName) {
		instance.HostName = localIp
	}
	if isEmpty(instance.VipAddress) {
		instance.VipAddress = strings.ToLower(config.InstanceConfig.AppName)
	}
	if isEmpty(instance.SecureVipAddress) {
		instance.SecureVipAddress = strings.ToLower(config.InstanceConfig.AppName)
	}

	//指示是否应在eureka注册后立即启用实例以获取流量,不建议立即开启
	if config.InstanceConfig.InstanceEnabledOnInit {
		instance.Status = core.STATUS_UP
	}

	//
	if isEmpty(instance.HomePageUrl) {
		instance.HomePageUrl = ""
	}
	if isEmpty(instance.StatusPageUrl) {
		instance.StatusPageUrl = "/actuator/info"
	}
	if isEmpty(instance.HealthCheckUrl) {
		instance.HealthCheckUrl = "/actuator/health"
	}
	instance.HomePageUrl = fmt.Sprintf("http://%s:%d/%s", localIp, port, strings.TrimLeft(instance.HomePageUrl, "/"))
	instance.StatusPageUrl = fmt.Sprintf("http://%s:%d/%s", localIp, port, strings.TrimLeft(instance.StatusPageUrl, "/"))
	instance.HealthCheckUrl = fmt.Sprintf("http://%s:%d/%s", localIp, port, strings.TrimLeft(instance.HealthCheckUrl, "/"))

	return instance, nil
}

func isEmpty(str string) bool {
	if 0 == len(str) {
		return true
	}
	str = strings.TrimSpace(str)
	if 0 == len(str) {
		return true
	}
	return false
}

// NewInstance 创建服务实例
func NewDefaultInstance() *core.Instance {
	port := 8080
	ip := netUtil.GetLocalIp()
	instance := &core.Instance{
		InstanceId:       ip + ":" + strconv.Itoa(port),
		HostName:         ip,
		App:              ip,
		IpAddr:           ip,
		VipAddress:       ip,
		SecureVipAddress: ip,
		Status:           core.STATUS_STARTING,
		Port: &core.Port{
			Port:    port,
			Enabled: "true",
		},
		SecurePort: &core.Port{
			Port:    443,
			Enabled: "false",
		},
		HomePageUrl:    "",
		StatusPageUrl:  "",
		HealthCheckUrl: "",
		// 数据中心
		DataCenterInfo: &core.DataCenterInfo{
			Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
			Name:  core.DC_NAME_TYPE_MY_OWN,
		},
		LeaseInfo: &core.LeaseInfo{
			RenewalIntervalInSecs: 30,
			DurationInSecs:        30,
		},
		CountryID: 0,
	}

	instance.HomePageUrl = fmt.Sprintf("http://%s:%d", ip, port)
	instance.StatusPageUrl = fmt.Sprintf("http://%s:%d%s", ip, port, "/actuator/info")
	instance.HealthCheckUrl = fmt.Sprintf("http://%s:%d%s", ip, port, "/actuator/health")
	return instance
}

//指示从eureka服务器获取注册表信息的频率,默认30秒
func (config *ClientConfig) getRegistryFetchIntervalSeconds() int {
	if 0 >= config.RegistryFetchIntervalSeconds {
		return 30
	}
	return config.RegistryFetchIntervalSeconds
}

func substitute(in []byte) ([]byte, error) {
	t, err := template.New("config").Parse(string(in))
	if err != nil {
		return nil, err
	}

	data := &templateData{
		Env: make(map[string]string),
	}

	values := os.Environ()
	for _, val := range values {
		keyval := strings.SplitN(val, "=", 2)
		if len(keyval) != 2 {
			continue
		}
		data.Env[keyval[0]] = keyval[1]
	}

	buffer := &bytes.Buffer{}
	if err = t.Execute(buffer, data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func validate(object interface{}) error {
	if valid := validateValue(object); valid != nil {
		return errors.New(fmt.Sprintf("Missing required config field: %v of type %s", valid.Name, valid.Type))
	}
	return nil
}

func validateValue(object interface{}) *inValid {
	objType := reflect.TypeOf(object)
	objValue := reflect.ValueOf(object)
	// If object is a nil interface value, TypeOf returns nil.
	if objType == nil {
		// Don't validate nil interfaces
		return nil
	}

	switch objType.Kind() {
	case reflect.Ptr:
		// If the ptr is nil
		if objValue.IsNil() {
			return &inValid{Type: objType.String()}
		}
		// De-reference the ptr and pass the object to validate
		return validateValue(objValue.Elem().Interface())
	case reflect.Struct:
		for idx := 0; idx < objValue.NumField(); idx++ {
			if valid := validateValue(objValue.Field(idx).Interface()); valid != nil {
				field := objType.Field(idx)
				// Capture sub struct names
				if valid.Name != "" {
					field.Name = field.Name + "." + valid.Name
				}

				// If our field is a pointer and it's pointing to an object
				if field.Type.Kind() == reflect.Ptr && !objValue.Field(idx).IsNil() {
					// The optional doesn't apply because our field does exist
					// instead the de-referenced object failed validation
					if field.Tag.Get("config") == "optional" {
						return &inValid{Name: field.Name, Type: valid.Type}
					}
				}
				// If the field is optional, don't invalidate
				if field.Tag.Get("config") != "optional" {
					return &inValid{Name: field.Name, Type: valid.Type}
				}
			}
		}
	// no way to tell if boolean or integer fields are provided or not
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Bool, reflect.Interface,
		reflect.Func:
		return nil
	default:
		if objValue.Len() == 0 {
			return &inValid{Type: objType.Name()}
		}
	}
	return nil
}
