//
// Copyright (c) 2017 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	hydro "github.com/openshift/ansible-service-broker/pkg/hydro/app"
	"github.com/openshift/ansible-service-broker/pkg/hydro/demobroker"
	"github.com/openshift/ansible-service-broker/pkg/hydro/log"
	"github.com/openshift/ansible-service-broker/pkg/hydro/server"
)

func main() {
	app := hydro.NewApp(demobroker.NewDemoBroker(), hydro.Config{
		log.LogConfig{Stdout: true, Level: "debug", Color: true},
		server.HandlerConfig{RequestDebug: true, Prefix: "/"},
	})
	app.Start()

	//app.RegisterBroker(NewAnsibleBroker())
	//app.RegisterJobs(NewJobManifest())
	//app.RegisterSubscribers(NewSubscriberManifest())
	//err := app.Start()
}
