package demobroker

import (
	"fmt"
	"github.com/openshift/ansible-service-broker/pkg/hydro/osb"
	"github.com/pborman/uuid"
)

type DemoBroker struct{}

func NewDemoBroker() *DemoBroker {
	return &DemoBroker{}
}

func (b *DemoBroker) Catalog() (*osb.CatalogResponse, error) {
	fmt.Println("DemoBroker::Catalog")
	return &osb.CatalogResponse{}, nil
}

func (b *DemoBroker) Provision(
	uuid.UUID, *osb.ProvisionRequest, bool,
) (*osb.ProvisionResponse, error) {
	fmt.Println("DemoBroker::Provision")
	return &osb.ProvisionResponse{}, nil
}

func (b *DemoBroker) Deprovision(
	osb.ServiceInstance, string, bool,
) (*osb.DeprovisionResponse, error) {
	fmt.Println("DemoBroker::Deprovision")
	return &osb.DeprovisionResponse{}, nil
}

func (b *DemoBroker) Bind(
	osb.ServiceInstance, uuid.UUID, *osb.BindRequest, bool,
) (*osb.BindResponse, bool, error) {
	fmt.Println("DemoBroker::Bind")
	return &osb.BindResponse{} /*was async?*/, true, nil
}

func (b *DemoBroker) Unbind(
	osb.ServiceInstance, osb.BindInstance, string, bool,
) (*osb.UnbindResponse, error) {
	fmt.Println("DemoBroker::Unbind")
	return &osb.UnbindResponse{}, nil
}

func (b *DemoBroker) Update(
	uuid.UUID, *osb.UpdateRequest, bool,
) (*osb.UpdateResponse, error) {
	fmt.Println("DemoBroker::Update")
	return &osb.UpdateResponse{}, nil
}

func (n *DemoBroker) LastOperation(
	uuid.UUID, *osb.LastOperationRequest,
) (*osb.LastOperationResponse, error) {
	fmt.Println("DemoBroker::LastOperation")
	return &osb.LastOperationResponse{}, nil
}

func (n *DemoBroker) GetServiceInstance(uuid.UUID) (*osb.ServiceInstance, error) {
	fmt.Println("DemoBroker::GetServiceInstance")
	return &osb.ServiceInstance{}, nil
}

func (n *DemoBroker) GetBindInstance(uuid.UUID) (*osb.BindInstance, error) {
	fmt.Println("DemoBroker::GetBindInstance")
	return &osb.BindInstance{}, nil
}
