package iregister

import (
	"math/rand"
	"strconv"
	"time"
)

type ServiceInstance interface {

	// GetInstanceId return The unique instance Id as registered.
	GetInstanceId() string

	// GetServiceId return The service Id as Registered.
	GetServiceId() string

	// GetHost return The hostname of the registered service instance.
	GetHost() string

	// GetPort return The port of the registered service instance.
	GetPort() int

	// IsSecure return Whether the port of the registered service instance.
	IsSecure() bool

	// GetMetadata return The key / value pair metadata associated with the service instance.
	GetMetadata() map[string]string
}

type DefaultServiceInstance struct {
	InstanceId string            // 具体的实例Id
	ServiceId  string            // 服务Id
	Host       string            // 主机地址
	Port       int               // 主机端口
	Secure     bool              // 是否安全
	Metadata   map[string]string // 元数据
}

func (d *DefaultServiceInstance) GetInstanceId() string {
	return d.InstanceId
}

func (d *DefaultServiceInstance) GetServiceId() string {
	return d.ServiceId
}

func (d *DefaultServiceInstance) GetHost() string {
	return d.Host
}

func (d *DefaultServiceInstance) GetPort() int {
	return d.Port
}

func (d *DefaultServiceInstance) IsSecure() bool {
	return d.Secure
}

func (d *DefaultServiceInstance) GetMetadata() map[string]string {
	return d.Metadata
}

func NewDefaultServiceInstance(serviceId string, host string, port int, secure bool, metadata map[string]string, instanceId string) (*DefaultServiceInstance, error) {

	if len(instanceId) == 0 {
		instanceId = serviceId + "-" + strconv.FormatInt(time.Now().Unix(), 10) + "-" + strconv.Itoa(rand.Intn(9000)+1000)
	}
	return &DefaultServiceInstance{InstanceId: instanceId, ServiceId: serviceId, Host: host, Port: port, Secure: secure, Metadata: metadata}, nil
}
