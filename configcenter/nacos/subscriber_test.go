package nacos

import (
	"fmt"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
}

func TestNacosSubscriber_Value(t *testing.T) {
	fmt.Println("Starting Nacos Subscriber Test...")
	nacosClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": []constant.ServerConfig{
			{
				IpAddr: "localhost",
				Port:   8848,
			},
		},
		"clientConfig": constant.ClientConfig{
			NamespaceId:         "test",
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "/tmp/nacos/log",
			CacheDir:            "/tmp/nacos/cache",
			Username:            "nacos",
			Password:            "nacos",
		},
	})

	if err != nil {
		t.Fatalf("failed to create nacos config client: %v", err)
	}

	subscriber := MustNacosSubscriber(nacosClient, "test", "demo", "DEFAULT_GROUP")

	center, err := configurator.NewConfigCenter[zrpc.RpcServerConf](configurator.Config{Type: "yaml", Log: true}, subscriber)
	if err != nil {
		t.Fatalf("failed to create config center: %v", err)
	}
	c, err := center.GetConfig()
	if err != nil {
		t.Fatalf("failed to get config: %v", err)
	}
	fmt.Printf("Config: %v\n", c)
	center.AddListener(func() {
		v, err := center.GetConfig()
		if err != nil {
			t.Fatalf("failed to get config: %v", err)
		}
		fmt.Printf("Config changed: %v\n", v)
	})
	fmt.Println("Waiting for config change...")
	time.Sleep(300 * time.Second)
}
