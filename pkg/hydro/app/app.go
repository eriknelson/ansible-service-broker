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

package app

import (
	"github.com/gorilla/handlers"
	"github.com/openshift/ansible-service-broker/pkg/hydro/osb"
	"github.com/openshift/ansible-service-broker/pkg/hydro/server"
	"os"
)

type HydroBroker interface {
	osb.OpenServiceBroker
	Initialize() error
}

type App struct {
	broker HydroBroker
	server server.Server
}

func NewApp(broker HydroBroker) App {
	return App{broker, server.NewDefaultServer()}
}

func NewAppWithServer(broker HydroBroker, server server.Server) App {
	return App{broker, server}
}

func (a *App) SetServer(server server.Server) {
	a.server = server
}

func (a *App) Start() {
	a.broker.Initialize()

	router := server.NewOpenServiceBrokerHandler(
		a.broker, server.HandlerConfig{true, a.server.Prefix()})
	if srv, ok := a.server.(server.RouterExtender); ok {
		srv.ExtendRouter(router)
	}
	handler := handlers.LoggingHandler(os.Stdout, router)
	a.server.StartServer(&handler)
}
