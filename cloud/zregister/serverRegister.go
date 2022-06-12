package zregister

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"strconv"
	"zinx_consul/cloud/iregister"
)

// consulServiceRegistry 服务注册实现
type consulServiceRegistry struct {
	serviceInstances     map[string]map[string]iregister.ServiceInstance
	client               *api.Client
	localServiceInstance iregister.ServiceInstance
}

// Register 注册服务实例到consul
func (c *consulServiceRegistry) Register(instance iregister.ServiceInstance) bool {
	// 创建注册到consul的服务到
	registration := new(api.AgentServiceRegistration) // Consul代理
	registration.ID = instance.GetInstanceId()
	registration.Name = instance.GetServiceId()
	registration.Port = instance.GetPort()
	registration.Address = instance.GetHost()
	var tags []string
	if instance.IsSecure() {
		tags = append(tags, "secure=true")
	} else {
		tags = append(tags, "secure=false")
	}

	if instance.GetMetadata() != nil {
		for key, val := range instance.GetMetadata() {
			tags = append(tags, key+"="+val)
		}
		registration.Tags = tags
	}
	registration.Tags = tags

	// 增加consul健康检查回调函数
	check := new(api.AgentServiceCheck)
	schema := "http"
	if instance.IsSecure() {
		schema = "https"
	}
	check.HTTP = fmt.Sprintf("%s://%s:%d/actuator/health", schema, registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "20s" //故障检查30s后， consul 自动将注册服务删除
	registration.Check = check

	// 注册服务到consul
	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		fmt.Println("register service to consul error:", err)
		return false
	}

	if c.serviceInstances == nil {
		c.serviceInstances = map[string]map[string]iregister.ServiceInstance{}
	}

	// 获取Service ，感觉这里的封装有问题
	services := c.serviceInstances[instance.GetServiceId()]

	c.serviceInstances[instance.GetInstanceId()] = services

	c.localServiceInstance = instance

	return true
}

// Deregister 注销服务
func (c *consulServiceRegistry) Deregister() {
	if c.serviceInstances == nil {
		return
	}

	// 获取本地服务实例列表
	services := c.serviceInstances[c.localServiceInstance.GetServiceId()]

	if services == nil {
		return
	}

	// 从本地列表中删除对应的服务实例
	delete(services, c.localServiceInstance.GetInstanceId())

	// 如果获取的服务实例列表为空，那么直接删除存储的 map[string]instance
	if len(services) == 0 {
		delete(c.serviceInstances, c.localServiceInstance.GetServiceId())
	}

	// 通知注册中心注销该服务
	err := c.client.Agent().ServiceDeregister(c.localServiceInstance.GetServiceId())
	if err != nil {
		fmt.Println("Deregister service error:", err)
		return
	}

	// 将结构体中的服务实例置空
	c.localServiceInstance = nil
}

// NewConsulServiceRegistry 外部接口，服务写完之后可以通过该函数进行注册
func NewConsulServiceRegistry(host string, port int, token string) (*consulServiceRegistry, error) {
	if len(host) < 3 {
		return nil, errors.New("check you host")
	}

	if port <= 0 || port > 65535 {
		return nil, errors.New("check port, port should between 1 and 65535")
	}

	config := api.DefaultConfig()
	config.Address = host + ":" + strconv.Itoa(port)
	config.Token = token
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &consulServiceRegistry{client: client}, nil
}
