package idiscovery

import "zinx_consul/cloud/iregister"

type DiscoveryClient interface {

	// GetInstances 获取所有的服务实例列表
	GetInstances(serviceId string) ([]iregister.ServiceInstance, error)

	// GetServices 获取所有的服务名称
	GetServices() ([]string, error)
}
