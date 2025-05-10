package bootstrap

import (
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
)

type NacosConfig struct {
	Addr      []string
	Username  string
	Password  string
	Namespace string
	DataId    string
	Group     string
}

type BootstrapConfig struct {
	Type  string      `json:"type"`
	Nacos NacosConfig `json:"nacos"`
}

func BootStrap(bootConfig BootstrapConfig, v any) config_client.IConfigClient {
	serverConfigs := make([]constant.ServerConfig, len(bootConfig.Nacos.Addr))
	for _, addr := range bootConfig.Nacos.Addr {
		parts := strings.Split(addr, ":")
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			IpAddr: parts[0],
			Port:   uint64(port),
		})
	}
	nacosClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig": constant.ClientConfig{
			NamespaceId:         bootConfig.Nacos.Namespace,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "/tmp/nacos/log",
			CacheDir:            "/tmp/nacos/cache",
			Username:            bootConfig.Nacos.Username,
			Password:            bootConfig.Nacos.Password,
		},
	})
	if err != nil {
		panic(err)
	}
	return nacosClient
}
