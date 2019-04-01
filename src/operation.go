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
	"fmt"
	"log"
	"time"

	"github.com/ystia/yorc/v3/config"
	"github.com/ystia/yorc/v3/events"
	"github.com/ystia/yorc/v3/prov"
)

type operationExecutor struct{}

func (d *operationExecutor) ExecAsyncOperation(ctx context.Context, conf config.Configuration, taskID, deploymentID, nodeName string, operation prov.Operation, stepName string) (*prov.Action, time.Duration, error) {
	return nil, 0, fmt.Errorf("*************asynchronous operations %v not yet supported by this sample", operation)
}

func (d *operationExecutor) ExecOperation(ctx context.Context, cfg config.Configuration, taskID, deploymentID, nodeName string, operation prov.Operation) error {

	log.Printf("************Hello from myOperationExecutor")
	_, err := cfg.GetConsulClient()
	if err != nil {
		return err
	}

	if cfg.Infrastructures["my-plugin"] != nil {
		for _, k := range cfg.Infrastructures["my-plugin"].Keys() {
			log.Printf("configuration key: %s", k)
		}
		log.Printf("Secret key: %q", cfg.Infrastructures["my-plugin"].GetStringOrDefault("test", "not found!"))
	}

	// Emit a log or an event
	events.WithContextOptionalFields(ctx).NewLogEntry(events.LogLevelINFO, deploymentID).Registerf("******Executing opeartion %q on node %q", operation.Name, nodeName)
	return nil
}
