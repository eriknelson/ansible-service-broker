//
// Copyright (c) 2018 Red Hat, Inc.
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
	//"encoding/json"
	"flag"
	"fmt"
	"net/http"
	//"os"

	//"github.com/automationbroker/bundle-lib/apb"
	crd "github.com/openshift/ansible-service-broker/pkg/dao/crd"

	"github.com/sirupsen/logrus"
)

var options struct {
	BrokerNamespace string
	Port            int
}

func init() {
	flag.IntVar(&options.Port, "port", 1337, "port that the dashboard-redirector should listen on")
	flag.StringVar(&options.BrokerNamespace, "ansible-service-broker", "", "namespace that the broker resides in")
	flag.Parse()
}

var crdDao *crd.Dao

var ServiceInstanceID = "XXX"

func main() {
	var err error
	logrus.Infof("Hello world.")

	crdDao, err = crd.NewDao(options.BrokerNamespace)
	if err != nil {
		panic(fmt.Sprintf("Unable to create crd client - %v", err))
	}

	http.HandleFunc("/", redirect)
	portStr := fmt.Sprintf(":%d", options.Port)

	logrus.Infof("Dashboard redirector listening on port [%s]", portStr)
	err = http.ListenAndServe(portStr, nil)
	if err != nil {
		logrus.Fatal("ListenAndServe: ", err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	specs, err := crdDao.BatchGetSpecs("")
	if err != nil {
		logrus.Errorf("Something went wrong trying to batch get specs! %+v", err.Error())
	} else {
		logrus.Infof("Got specs! %+v", specs)
	}

	http.Redirect(w, r, "http://www.google.com", 301)
}
