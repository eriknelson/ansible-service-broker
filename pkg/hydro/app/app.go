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
	"github.com/openshift/ansible-service-broker/pkg/hydro/log"
	"github.com/openshift/ansible-service-broker/pkg/hydro/osb"
	"github.com/openshift/ansible-service-broker/pkg/hydro/server"
	"os"
)

type App struct {
	broker osb.OpenServiceBroker
	server server.Server
	config Config
}

type Config struct {
	LogConfig     log.LogConfig
	HandlerConfig server.HandlerConfig
}

func NewApp(broker osb.OpenServiceBroker, config Config) App {
	log.InitializeLog(config.LogConfig)
	return App{broker, server.NewDefaultServer(), config}
}

func (a *App) SetServer(server server.Server) {
	a.server = server
}

func (a *App) Start() {
	router := server.NewOpenServiceBrokerHandler(a.broker, a.config.HandlerConfig)
	if srv, ok := a.server.(server.RouterExtender); ok {
		srv.ExtendRouter(router)
	}
	handler := handlers.LoggingHandler(os.Stdout, router)
	a.server.StartServer(&handler)
}
