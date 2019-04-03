// Copyright 2018 Bull S.A.S. Atos Technologies - Bull, Rue Jean Jaures, B.P.68, 78340, Les Clayes-sous-Bois, France.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"

	"github.com/ystia/yorc/v3/plugin"
	"github.com/ystia/yorc/v3/prov"

)

func main() {
	// Create configuration that defines the type of plugins to be served.
	// In servConfig can be set :
	// - TOSCA definitions for an extended Yorc
	// - A DelegateExecutor for some TOSCA component types
	// - An OperationExecutor for some TOSCA atrifacts types
	// - An InfrastructureUsageCollector for specific instrastructures to be monitored
	var servConfig *plugin.ServeOpts
	servConfig = new(plugin.ServeOpts)

	// Add TOSCA Definitions contained in the def variable.
	// For this example, TOSCA definitions are provided in a yaml file mycustom-types.yaml
	// bundled in this binary (see Makefile)
	// The mycustom-types.yaml key can be then by used (imported) by applications
	// deployed to the extended Yorc, as the example tosca/topology.yaml is doing
	var err error
	servConfig.Definitions, err = getToscaResources()
	if err != nil {
		log.Printf("Error getting bundle TOSCA resources: %+v\n", err)
		return
	}

	// Set DelegateFunc that implements a DelegateExecutor for the TOSCA component types specified in DelegateSupportedTypes
	// The delegateExecutor is defined in delegate.go
	servConfig.DelegateSupportedTypes = []string{`mytosca\.types\..*`}
	servConfig.DelegateFunc = func() prov.DelegateExecutor {
		return new(delegateExecutor)
	}

	// Set OperationFunc that implements an OperationExecutor for the TOSCA artifacts specified in OperationSupportedArtifactTypes
	servConfig.OperationSupportedArtifactTypes = []string{"mytosca.artifacts.Implementation.MyImplementation"}
	servConfig.OperationFunc = func() prov.OperationExecutor {
		return new(operationExecutor)
	}

	plugin.Serve(servConfig)
}