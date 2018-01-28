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

package server

import (
	"github.com/openshift/ansible-service-broker/pkg/hydro/log"
	"net/http"
	"time"
)

type DefaultServer struct {
	prefix string
}

func NewDefaultServer() *DefaultServer {
	return &DefaultServer{prefix: "/"}
}

func (s *DefaultServer) SetPrefix(prefix string) {
	s.prefix = prefix
}

func (s *DefaultServer) Prefix() string {
	return s.prefix
}

func (s *DefaultServer) StartServer(h *http.Handler) {
	addr := "0.0.0.0:1338"
	log.Infof("DefaultServer Listening on %s", addr)

	srv := &http.Server{
		Handler:      *h,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
}
