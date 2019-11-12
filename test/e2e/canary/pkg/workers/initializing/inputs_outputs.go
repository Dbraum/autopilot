package initializing

import (
	parameters "github.com/solo-io/autopilot/test/e2e/canary/pkg/parameters"
)

type Inputs struct {
	Deployments parameters.Deployments
}

// FindDeployment returns <Deployment, true> if the item is found. else parameters.Deployment{}, false
func (i Inputs) FindDeployment(name, namespace string) (parameters.Deployment, bool) {
	for _, item := range i.Deployments.Items {
		if item.Name == name && item.Namespace == namespace {
			return item, true
		}
	}
	return parameters.Deployment{}, false
}

type Outputs struct {
	Deployments     parameters.Deployments
	Services        parameters.Services
	VirtualServices parameters.VirtualServices
}
