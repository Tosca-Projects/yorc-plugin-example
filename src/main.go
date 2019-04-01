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
	"github.com/ystia/yorc/v3/plugin"
	"github.com/ystia/yorc/v3/prov"
)

// TOSCA definition content
// (a yaml file content)
var def = []byte(`tosca_definitions_version: yorc_tosca_simple_yaml_1_0

metadata:
  template_name: mytosca-types
  template_version: 1.0.0
  template_author: yorc

imports:
  - <normative-types.yml>

artifact_types:
  mytosca.artifacts.Implementation.MyImplementation:
    derived_from: tosca.artifacts.Implementation
    description: My dummy implementation artifact
    file_ext: [ myext ]

node_types:
  mytosca.types.Compute:
    derived_from: tosca.nodes.Compute

  mytosca.types.Soft:
    derived_from: tosca.nodes.SoftwareComponent
    interfaces:
      Standard:
        create: dothis.myext
`)

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
	// The mycustom-types.yaml key can be used (imported) by applications deployed to the extended Yorc
	servConfig.Definitions = map[string][]byte{"mycustom-types.yaml": def}

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
