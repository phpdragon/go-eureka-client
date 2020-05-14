package core

import (
	"fmt"
	httpClient "github.com/phpdragon/go-eurake-client/httpclent"
	"net/url"
	"strings"
)

// wiki: https://github.com/Netflix/eureka/wiki/Eureka-REST-operations
type EurekaServerApi struct {
	BaseUrl string
}

func NewEurekaServerApi(baseUrl string) *EurekaServerApi {
	return &EurekaServerApi{
		BaseUrl: baseUrl,
	}
}

func (api *EurekaServerApi) url(path string) string {
	return strings.TrimRight(api.BaseUrl, "/") + path
}

// Register new application instance by brief info
func (api *EurekaServerApi) RegisterInstance(appId string, instance *Instance) error {
	eurekaUrl := api.url("/apps/" + strings.ToUpper(appId))

	body := map[string]interface{}{"instance": instance}
	// status: httpClient.StatusNoContent
	result := httpClient.Post(eurekaUrl).Json(body).Send().Status2xx()

	if result.Err != nil {
		return fmt.Errorf("Register application instance failed, error: %s", result.Err)
	}

	return nil
}

// 更新实例状态
// update instance status
func (api *EurekaServerApi) UpdateInstanceStatus(appId, instanceId, status string) error {
	eurekaUrl := api.url(fmt.Sprintf("/apps/%s/%s/status?value=%s", strings.ToUpper(appId), instanceId, status))
	// status: httpClient.StatusNoContent
	result := httpClient.Put(eurekaUrl).Send().StatusOk()

	if result.Err != nil {
		return fmt.Errorf("ClientConfig UP failed, err=%s", result.Err)
	}

	return nil
}

// 更新实例的元数据
// Update metadata
func (api *EurekaServerApi) UpdateMeta(appId, instanceId string, metadata map[string]string) error {
	queryStr := ""
	for k, v := range metadata {
		queryStr += fmt.Sprintf("&%s=%s", k, v)
	}
	queryStr = strings.TrimLeft(queryStr, "&")

	eurekaUrl := api.url(fmt.Sprintf("/apps/%s/%s/metadata?%s", appId, instanceId, queryStr))
	// status: httpClient.StatusNoContent
	result := httpClient.Put(eurekaUrl).Send().StatusOk()
	if result.Err != nil {
		return fmt.Errorf("Failed to update instance metadata, err=%s", result.Err)
	}

	return nil
}

// Heartbeat 发送心跳
// PUT /eureka/v2/apps/appID/instanceID
func (api *EurekaServerApi) SendHeartbeat(appId, instanceID string) error {
	eurekaUrl := api.url("/apps/" + strings.ToUpper(appId) + "/" + instanceID)
	params := url.Values{
		"status": {"UP"},
	}

	result := httpClient.Put(eurekaUrl).Params(params).Send().StatusOk()
	if result.Err != nil {
		return fmt.Errorf("Heartbeat failed, error: %s", result.Err)
	}
	return nil
}

// DeRegisterInstance 删除实例
// DELETE /eureka/v2/apps/appID/instanceID
func (api *EurekaServerApi) DeRegisterInstance(appId, instanceID string) error {
	eurekaUrl := api.url("/apps/" + strings.ToUpper(appId) + "/" + instanceID)

	// status: httpClient.StatusNoContent
	result := httpClient.Delete(eurekaUrl).Send().StatusOk()
	if result.Err != nil {
		return fmt.Errorf("UnRegister application instance failed, error: %s", result.Err)
	}
	return nil
}

// Refresh 查询所有服务实例
// GET /eureka/v2/apps
// Query for all instances
func (api *EurekaServerApi) QueryAllInstances() (*Applications, error) {
	eurekaUrl := api.url("/apps")
	res := &EurekaApps{}

	err := httpClient.Get(eurekaUrl).Header("Accept", " application/json").Send().StatusOk().Json(&res)
	if err != nil {
		return nil, fmt.Errorf("Refresh failed, error: %s", err.Error())
	}
	return &res.Applications, nil
}

//GET /eureka/v2/apps/appID
// Query for all instances by appId
func (api *EurekaServerApi) QueryAllInstanceByAppId(appId string) (*Application, error) {
	eurekaUrl := api.url("/apps/" + strings.ToUpper(appId))
	res := &EurekaApp{}
	err := httpClient.Get(eurekaUrl).Header("Accept", " application/json").Send().StatusOk().Json(&res)
	if err != nil {
		return nil, fmt.Errorf("Failed to query appId instances, err=%s", err.Error())
	}
	return &res.Application, nil
}

// 查询单个实例详情
// query specific instanceId
func (api *EurekaServerApi) QuerySpecificAppInstance(instanceId string) (*Instance, error) {
	eurekaUrl := api.url("/instances/" + instanceId)
	res := &EurekaInstance{}
	err := httpClient.Get(eurekaUrl).Header("Accept", " application/json").Send().StatusOk().Json(&res)
	if err != nil {
		return nil, fmt.Errorf("Failed to query appId instances, err=%s", err.Error())
	}
	return &res.Instance, nil
}

//Query for all instances under a particular vip address
func (api *EurekaServerApi) QueryAllInstancesByVipAddress(vipAddress string) (*Applications, error) {
	eurekaUrl := api.url("/vips/" + vipAddress)
	res := &EurekaApps{}

	err := httpClient.Get(eurekaUrl).Header("Accept", " application/json").Send().StatusOk().Json(&res)
	if err != nil {
		return nil, fmt.Errorf("Failed to query appId instances, err=%s", err.Error())
	}
	return &res.Applications, nil
}

//Query for all instances under a particular secure vip address
func (api *EurekaServerApi) QueryAllInstancesBySvipAddress(svipAddress string) (*Applications, error) {
	eurekaUrl := api.url("/svips/" + svipAddress)
	res := &EurekaApps{}
	err := httpClient.Get(eurekaUrl).Header("Accept", " application/json").Send().StatusOk().Json(&res)
	if err != nil {
		return nil, fmt.Errorf("Failed to query appId instances, err=%s", err.Error())
	}
	return &res.Applications, nil
}
