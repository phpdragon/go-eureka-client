package core

const (
	STATUS_UP             = "UP"
	STATUS_DOWN           = "DOWN"
	STATUS_STARTING       = "STARTING"
	STATUS_OUT_OF_SERVICE = "OUT_OF_SERVICE"
	STATUS_UNKNOWN        = "UNKNOWN"

	DC_NAME_TYPE_MY_OWN = "MyOwn"
	DC_NAME_TYPE_AMAZON = "Amazon"
)

type (
	EurekaApps struct {
		Applications Applications `json:"applications,omitempty"`
	}

	EurekaApp struct {
		Application Application `json:"application,omitempty"`
	}

	EurekaInstance struct {
		Instance Instance `json:"instance,omitempty"`
	}

	// applications eureka服务端注册的apps
	Applications struct {
		VersionsDelta string        `json:"versions__delta,omitempty"`
		AppsHashcode  string        `json:"apps__hashcode,omitempty"`
		Applications  []Application `json:"application,omitempty"`
	}

	// Application eureka服务端注册的app
	Application struct {
		Name      string     `json:"name"`
		Instances []Instance `json:"instance"`
	}

	// InstanceConfig 服务实例
	Instance struct {
		// Register application instance needed -- BEGIN
		InstanceId       string          `json:"instanceId,omitempty"`
		HostName         string          `json:"hostName"`
		App              string          `json:"app"`
		IpAddr           string          `json:"ipAddr"`
		Status           string          `json:"status"`
		VipAddress       string          `json:"vipAddress"`
		SecureVipAddress string          `json:"secureVipAddress,omitempty"`
		Port             *Port           `json:"port,omitempty"`
		SecurePort       *Port           `json:"securePort,omitempty"`
		HomePageUrl      string          `json:"homePageUrl,omitempty"`
		StatusPageUrl    string          `json:"statusPageUrl"`
		HealthCheckUrl   string          `json:"healthCheckUrl,omitempty"`
		DataCenterInfo   *DataCenterInfo `json:"dataCenterInfo"`
		LeaseInfo        *LeaseInfo      `json:"leaseInfo,omitempty"`
		// Register application instance needed -- END

		OverriddenStatus              string                 `json:"overriddenstatus,omitempty"`
		LastUpdatedTimestamp          string                 `json:"lastUpdatedTimestamp,omitempty"`
		LastDirtyTimestamp            string                 `json:"lastDirtyTimestamp,omitempty"`
		ActionType                    string                 `json:"actionType,omitempty"`
		Metadata                      map[string]interface{} `json:"metadata,omitempty"`
		IsCoordinatingDiscoveryServer string                 `json:"isCoordinatingDiscoveryServer,omitempty"`
		CountryID                     int                    `json:"countryId,omitempty"`
	}

	// Port 端口
	Port struct {
		Port    int    `json:"$"`
		Enabled string `json:"@enabled"`
	}

	// DataCenterInfo 数据中心信息
	DataCenterInfo struct {
		Name     string              `json:"name"`
		Class    string              `json:"@class"`
		Metadata *DataCenterMetadata `json:"metadata,omitempty"`
	}

	// DataCenterMetadata 数据中心信息元数据
	DataCenterMetadata struct {
		AmiLaunchIndex   string `json:"ami-launch-index,omitempty"`
		LocalHostname    string `json:"local-hostname,omitempty"`
		AvailabilityZone string `json:"availability-zone,omitempty"`
		InstanceId       string `json:"instance-id,omitempty"`
		PublicIpv4       string `json:"public-ipv4,omitempty"`
		PublicHostname   string `json:"public-hostname,omitempty"`
		AmiManifestPath  string `json:"ami-manifest-path,omitempty"`
		LocalIpv4        string `json:"local-ipv4,omitempty"`
		Hostname         string `json:"hostname,omitempty"`
		AmiID            string `json:"ami-id,omitempty"`
		InstanceType     string `json:"instance-type,omitempty"`
	}

	// LeaseInfo 续约信息
	LeaseInfo struct {
		RenewalIntervalInSecs int   `json:"renewalIntervalInSecs,omitempty"`
		DurationInSecs        int   `json:"durationInSecs,omitempty"`
		RegistrationTimestamp int64 `json:"registrationTimestamp,omitempty"`
		LastRenewalTimestamp  int64 `json:"lastRenewalTimestamp,omitempty"`
		EvictionTimestamp     int64 `json:"evictionTimestamp,omitempty"`
		ServiceUpTimestamp    int64 `json:"serviceUpTimestamp,omitempty"`
	}
)
