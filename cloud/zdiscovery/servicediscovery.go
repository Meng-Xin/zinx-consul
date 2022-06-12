package zdiscovery

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"strconv"
	"unsafe"
	"zinx_consul/cloud/iregister"
)

type consulServiceRegistry struct {
	serviceInstances     map[string]map[string]iregister.ServiceInstance
	client               *api.Client
	localServiceInstance iregister.ServiceInstance
}

func (c *consulServiceRegistry) GetInstances(serviceId string) ([]iregister.ServiceInstance, error) {
	service, _, err := c.client.Catalog().Service(serviceId, "", nil)
	if err != nil {
		return nil, err
	}
	if len(service) > 0 {
		result := make([]iregister.ServiceInstance, len(service))
		for idx, server := range service {
			s := iregister.DefaultServiceInstance{
				InstanceId: server.ServiceID,
				ServiceId:  server.ServiceName,
				Host:       server.Address,
				Port:       server.ServicePort,
				Metadata:   server.ServiceMeta,
			}
			result[idx] = &s
		}
		return result, nil
	}
	return nil, nil
}

func (c *consulServiceRegistry) GetServices() ([]string, error) {
	services, _, err := c.client.Catalog().Services(nil)
	if err != nil {
		return nil, err
	}
	result := make([]string, unsafe.Sizeof(services))
	index := 0
	for serviceName, _ := range services {
		result[index] = serviceName
		index++
	}
	return result, nil
}

// NewConsulServiceRegistry New a consulServiceRegistry instance
func NewConsulServiceRegistry(host string, port int, token string) (*consulServiceRegistry, error) {
	if len(host) < 3 {
		return nil, errors.New("check you host")
	}

	if port <= 0 || port > 65535 {
		return nil, errors.New("check port,port should between 1 and 65535")
	}

	config := api.DefaultConfig()
	config.Address = host + ":" + strconv.Itoa(port)
	config.Token = token
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Println("New Client fail error:", err)
		return nil, err
	}
	return &consulServiceRegistry{client: client}, nil

}
