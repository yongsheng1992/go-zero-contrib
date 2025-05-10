package nacos

import (
	"sync"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
)

type nacosSubscriber struct {
	listeners    []func()
	lock         sync.Mutex
	currentValue string
	nacosClient  config_client.IConfigClient
	DataId       string
	Group        string
	Namespace    string
}

func (s *nacosSubscriber) AddListener(listener func()) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.listeners = append(s.listeners, listener)
	return nil
}

func (s *nacosSubscriber) Value() (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.currentValue, nil
}

func (s *nacosSubscriber) onChange(namespace, group, dataId, data string) {
	if s.Namespace != namespace {
		return
	}
	if s.DataId != dataId {
		return
	}
	if s.Group != group {
		return
	}

	s.lock.Lock()
	s.currentValue = data
	s.lock.Unlock()

	for _, listener := range s.listeners {
		listener()
	}
}

func MustNacosSubscriber(nacosClient config_client.IConfigClient, namespace, group, dataId string) subscriber.Subscriber {
	configContent, err := nacosClient.GetConfig(vo.ConfigParam{DataId: dataId, Group: group})
	if err != nil {
		panic(err)
	}

	nacosSubscriber := &nacosSubscriber{
		listeners:    []func(){},
		currentValue: configContent,
		nacosClient:  nacosClient,
		Namespace:    namespace,
		DataId:       dataId,
		Group:        group,
	}

	if err := nacosClient.ListenConfig(vo.ConfigParam{DataId: dataId, Group: group, OnChange: nacosSubscriber.onChange}); err != nil {
		panic(err)
	}

	return nacosSubscriber
}

func MustNacosSubscriberV1(serverConfig []constant.ServerConfig, clientConfig constant.ClientConfig, namespace, group, dataId string) subscriber.Subscriber {
	nacosClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfig,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	return MustNacosSubscriber(nacosClient, namespace, group, dataId)
}
