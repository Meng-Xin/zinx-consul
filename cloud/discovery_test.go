package cloud

import (
	"testing"
	"zinx_consul/cloud/zdiscovery"
)

func TestConsulServiceDiscover(t *testing.T) {
	host := "127.0.0.1"
	port := 8500
	token := ""
	registry, err := zdiscovery.NewConsulServiceRegistry(host, port, token)
	if err != nil {
		return
	}

	t.Log(registry.GetServices())

	t.Log(registry.GetInstances("go-user-server"))
}
