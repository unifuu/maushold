package service

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

type ServiceDiscovery struct {
	consul *consulapi.Client
}

func NewServiceDiscovery(consul *consulapi.Client) *ServiceDiscovery {
	return &ServiceDiscovery{consul: consul}
}

// DiscoverService discovers a healthy service instance from Consul
func (sd *ServiceDiscovery) DiscoverService(serviceName string) (string, error) {
	if sd.consul == nil {
		return "", fmt.Errorf("consul client is not initialized")
	}

	// Query Consul for the player service
	services, _, err := sd.consul.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("service %s not found in Consul", serviceName)
	}

	if len(services) == 0 {
		return "", fmt.Errorf("service %s not found in Consul", serviceName)
	}

	// Use the first healthy service instance
	service := services[0]
	address := service.Service.Address
	port := service.Service.Port

	url := fmt.Sprintf("http://%s:%d", address, port)
	log.Printf("Discovered service %s at %s", serviceName, url)

	return url, nil
}

// DiscoverAllInstances discovers all healthy service instances from Consul
func (sd *ServiceDiscovery) DiscoverAllInstances(serviceName string) ([]string, error) {
	if sd.consul == nil {
		return nil, fmt.Errorf("consul client is not initialized")
	}

	// Query Consul for the player service
	services, _, err := sd.consul.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("service %s not found in Consul", serviceName)
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("service %s not found in Consul", serviceName)
	}

	urls := make([]string, len(services))
	for _, service := range services {
		address := service.Service.Address
		if address == "" {
			address = service.Node.Address
		}
		port := service.Service.Port
		urls = append(urls, fmt.Sprintf("http://%s:%d", address, port))
	}

	return urls, nil
}
