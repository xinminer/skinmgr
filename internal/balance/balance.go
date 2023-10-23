package balance

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"math/rand"
)

type DiscoveryConfig struct {
	ID      string
	Name    string
	Tags    []string
	Port    int
	Address string
	Meta    map[string]string
}

func RegisterService(addr string, dis DiscoveryConfig) error {
	config := consulapi.DefaultConfig()
	config.Address = addr
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Printf("create consul client : %v\n", err.Error())
		return err
	}
	registration := &consulapi.AgentServiceRegistration{
		ID:      dis.ID,
		Name:    dis.Name,
		Port:    dis.Port,
		Tags:    dis.Tags,
		Address: dis.Address,
		Meta:    dis.Meta,
	}

	check := &consulapi.AgentServiceCheck{}
	check.TCP = fmt.Sprintf("%s:%d", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "60s"
	registration.Check = check

	if err = client.Agent().ServiceRegister(registration); err != nil {
		return err
	}

	return nil
}

func Random(consulAddr, service, preferred string) (string, error) {
	target, err := getConsulServices(consulAddr, service)
	if err != nil {
		return "", err
	}
	lens := len(target)
	index := rand.Intn(lens)
	for i := 0; i < 3; i++ {
		svr := target[index]
		if svr.Service.Meta["preferred"] == preferred {
			return svr.Service.Meta["target"], nil
		}
		index = rand.Intn(lens)
	}
	svr := target[index]
	return svr.Service.Meta["target"], nil
}

func getConsulServices(consulAddr, service string) ([]*consulapi.ServiceEntry, error) {
	config := consulapi.DefaultConfig()
	config.Address = consulAddr
	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	services, _, err := client.Health().Service(service, "", false, nil)
	return services, err
}
