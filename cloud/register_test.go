package cloud

import (
	"github.com/gin-gonic/gin"
	"testing"
	"zinx_consul/cloud/iregister"
	"zinx_consul/cloud/zregister"
)

func TestConsulServiceRegistry(t *testing.T) {
	host := "127.0.0.1"
	port := 8500
	registryDiscoveryClient, _ := zregister.NewConsulServiceRegistry(host, port, "")

	// 创建服务实例
	instance, err := iregister.NewDefaultServiceInstance("go-user-server", "", 8090, false, map[string]string{"user": "zyn"}, "")
	if err != nil {
		return
	}
	registryDiscoveryClient.Register(instance)

	r := gin.Default()

	// 健康检测接口，只要是200就是成功
	r.GET("/actuator/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err = r.Run(":8090")
	if err != nil {
		registryDiscoveryClient.Deregister()
	}

}
