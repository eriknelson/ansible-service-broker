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
	"os"

	"github.com/automationbroker/bundle-lib/apb"
	crd "github.com/openshift/ansible-service-broker/pkg/dao/crd"

	"github.com/sirupsen/logrus"
)

var options struct {
	BrokerNamespace string
}

func init() {
	flag.StringVar(&options.BrokerNamespace, "namespace", "", "namespace that the broker resides in")
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

	// convert all the service instances
	siSaved := []*apb.ServiceInstance{}

	//if siJSONStrs != nil {
	//for _, str := range *siJSONStrs {
	//si := apb.ServiceInstance{}
	//err := json.Unmarshal([]byte(str), &si)
	//if err != nil {
	//revertServiceInstance(siSaved)
	//revertCrdSavedSpecs(crdSavedSpecs)
	//panic(fmt.Sprintf("Unable to migrate all the service instances json unmarshal error - %v", err))
	//}
	//err = crdDao.SetServiceInstance(si.ID.String(), &si)
	//if err != nil {
	//revertServiceInstance(siSaved)
	//revertCrdSavedSpecs(crdSavedSpecs)
	//panic(fmt.Sprintf("Unable to migrate all the service instances set service instance - %v", err))
	//}
	//siSaved = append(siSaved, &si)
	//}
	//}
}
