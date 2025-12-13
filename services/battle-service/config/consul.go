package config

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

func InitConsul(cfg *Config) *consulapi.Client {
	config := consulapi.DefaultConfig()
	config.Address = cfg.ConsulAddr

	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Printf("Failed to connect to Consul: %v", err)
		return nil
	}

	log.Println("Consul connected")
	return client
}

func RegisterService(client *consulapi.Client, serviceName, port string) error {
	if client == nil {
		return fmt.Errorf("consul client is nil")
	}

	registration := &consulapi.AgentServiceRegistration{
		ID:      serviceName,
		Name:    serviceName,
		Port:    parsePort(port),
		Address: serviceName, // Use service name as address in Docker network
		Check: &consulapi.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%s/health", serviceName, port),
			Interval:                       "10s",
			Timeout:                        "3s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	return client.Agent().ServiceRegister(registration)
}

func DeregisterService(client *consulapi.Client, serviceName string) {
	if client != nil {
		client.Agent().ServiceDeregister(serviceName)
	}
}

func parsePort(portStr string) int {
	var port int
	fmt.Sscanf(portStr, "%d", &port)
	return port
}
