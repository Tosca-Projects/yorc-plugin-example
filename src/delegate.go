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
	"context"

	"github.com/ystia/yorc/v4/config"
	"github.com/ystia/yorc/v4/deployments"
	"github.com/ystia/yorc/v4/events"
	"github.com/ystia/yorc/v4/locations"
	"github.com/ystia/yorc/v4/log"
	"github.com/ystia/yorc/v4/tasks"
	"github.com/ystia/yorc/v4/tosca"
)

type delegateExecutor struct{}

func (de *delegateExecutor) ExecDelegate(ctx context.Context, conf config.Configuration, taskID, deploymentID, nodeName, delegateOperation string) error {
	log.Debugf("Entering plugin ExecDelegate")

	// Here is how to retrieve location properties

	var locationProps config.DynamicMap
	locationMgr, err := locations.GetManager(conf)
	if err != nil {
		return err
	}

	locationProps, err = locationMgr.GetLocationProperties("my-location", "my-infra")
	if err != nil {
		return err
	}

	log.Debugf("********Got my-location properties")
	for k, v := range locationProps {
		events.WithContextOptionalFields(ctx).NewLogEntry(events.LogLevelINFO, deploymentID).Registerf("**********location property key: %q", k)
		events.WithContextOptionalFields(ctx).NewLogEntry(events.LogLevelINFO, deploymentID).Registerf("**********location property value: %q", v)
	}

	// TODO: add here the code retrieving properties to connect to the API
	// allowing to allocated compute instances/connect to your
	// infrastructure

	// Get node instances related to this task (may be a subset of all instances for a scaling operation for instance)
	instances, err := tasks.GetInstances(ctx, taskID, deploymentID, nodeName)
	if err != nil {
		return err
	}

	// Emit events and logs on instance status change
	for _, instanceName := range instances {
		deployments.SetInstanceStateWithContextualLogs(ctx, deploymentID, nodeName, instanceName, tosca.NodeStateCreating)
	}

	// Use the deployments api to get info about the node to provision
	nodeType, err := deployments.GetNodeType(ctx, deploymentID, nodeName)
	if err != nil {
		return err
	}

	// Emit a log or an event
	events.WithContextOptionalFields(ctx).NewLogEntry(events.LogLevelINFO, deploymentID).Registerf("**********Provisioning node %q of type %q", nodeName, nodeType)

	for _, instanceName := range instances {
		// TODO: add here the code allowing to create a Compute Instance
		deployments.SetInstanceStateWithContextualLogs(ctx, deploymentID, nodeName, instanceName, tosca.NodeStateStarted)
	}
	return nil
}
