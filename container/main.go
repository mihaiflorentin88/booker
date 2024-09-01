package container

var Load *ServiceContainer

type ServiceContainer struct {
	infrastructure *InfrastructureContainer
}

func NewServiceContainer() *ServiceContainer {
	return &ServiceContainer{
		infrastructure: &InfrastructureContainer{},
	}
}
