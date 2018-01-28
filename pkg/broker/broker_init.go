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

package broker

import (
	"context"
	"os"
	apirbac "k8s.io/api/rbac/v1beta1"
	"k8s.io/kubernetes/pkg/apis/rbac"
	"github.com/openshift/ansible-service-broker/pkg/clients"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	flags "github.com/jessevdk/go-flags"
	"github.com/openshift/ansible-service-broker/pkg/version"
	"fmt"
	"github.com/openshift/ansible-service-broker/pkg/config"
	agnosticruntime "github.com/openshift/ansible-service-broker/pkg/runtime"
	logutil "github.com/openshift/ansible-service-broker/pkg/util/logging"
	"github.com/openshift/ansible-service-broker/pkg/dao"
	"github.com/openshift/ansible-service-broker/pkg/registries"
	"github.com/openshift/ansible-service-broker/pkg/apb"
)

const (
	// msgBufferSize - The buffer for the message channel.
	msgBufferSize = 20
)

// args - Command line arguments for the ansible service broker.
type args struct {
	ConfigFile string `short:"c" long:"config" description:"Config File" default:"/etc/ansible-service-broker/config.yaml"`
	Version    bool   `short:"v" long:"version" description:"Print version information"`
}

// createArgs - Will return the arguments that were passed in to the application
func createArgs() (args, error) {
	args := args{}

	_, err := flags.Parse(&args)
	if err != nil {
		return args, err
	}
	return args, nil
}

func (a *AnsibleBroker) Initialize() error {
	var err error

	fmt.Println("============================================================")
	fmt.Println("==           Starting Ansible Broker...                   ==")
	fmt.Println("============================================================")


	// Load cli args
	if a.args, err = createArgs(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	// TODO: Args, config, versioning, and loggging are probably the domain of hydro since they're
	// generically applicable to brokers. Probably want to move these there in a strategic way

	if a.args.Version {
		fmt.Println(version.Version)
		os.Exit(0)
	}

	// Load config
	if a.config, err = config.CreateConfig(a.args.ConfigFile); err != nil {
		os.Stderr.WriteString("ERROR: Failed to read config file\n")
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	// Init logs
	c := logutil.LogConfig{
		LogFile: a.config.GetString("log.logfile"),
		Stdout:  a.config.GetBool("log.stdout"),
		Level:   a.config.GetString("log.level"),
		Color:   a.config.GetBool("log.color"),
	}
	if err = logutil.InitializeLog(c); err != nil {
		os.Stderr.WriteString("ERROR: Failed to initialize logger\n")
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	// TODO: Should really start to put these behind some functions for better encapsulation
	// Initializing clients as soon as we have deps ready.
	brokerConfig := a.config.GetSubConfig("broker")
	a.brokerConfig = Config{
		DevBroker:          brokerConfig.GetBool("dev_broker"),
		LaunchApbOnBind:    brokerConfig.GetBool("launch_apb_on_bind"),
		BootstrapOnStartup: brokerConfig.GetBool("bootstrap_on_startup"),
		Recovery:           brokerConfig.GetBool("recovery"),
		OutputRequest:      brokerConfig.GetBool("output_request"),
		SSLCertKey:         brokerConfig.GetString("ssl_cert_key"),
		SSLCert:            brokerConfig.GetString("ssl_cert"),
		RefreshInterval:    brokerConfig.GetString("refresh_interval"),
		AutoEscalate:       brokerConfig.GetBool("auto_escalate"),
		ClusterURL:         brokerConfig.GetString("cluster_url"),
	}

	err = initClients(a.config)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// Initialize Runtime
	log.Debug("Connecting to Cluster")
	agnosticruntime.NewRuntime()
	agnosticruntime.Provider.ValidateRuntime()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	if a.config.GetBool("broker.recovery") {
		log.Info("Initiating Recovery Process")
		a.Recover()
	}

	log.Debug("Connecting Dao")
	a.dao, err = dao.NewDao()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	log.Debug("Connecting Registry")
	for _, config := range a.config.GetSubConfigArray("registry") {
		reg, err := registries.NewRegistry(config, a.config.GetString("openshift.namespace"))
		if err != nil {
			log.Errorf(
				"Failed to initialize %v Registry err - %v \n", config.GetString("name"), err)
			os.Exit(1)
		}
		a.registry = append(a.registry, reg)
	}

	validateRegistryNames(a.registry)

	log.Debug("Initializing WorkEngine")
	a.engine = NewWorkEngine(msgBufferSize)
	err = a.engine.AttachSubscriber(
		NewProvisionWorkSubscriber(a.dao),
		ProvisionTopic)
	if err != nil {
		log.Errorf("Failed to attach subscriber to WorkEngine: %s", err.Error())
		os.Exit(1)
	}
	err = a.engine.AttachSubscriber(
		NewDeprovisionWorkSubscriber(a.dao),
		DeprovisionTopic)
	if err != nil {
		log.Errorf("Failed to attach subscriber to WorkEngine: %s", err.Error())
		os.Exit(1)
	}
	err = a.engine.AttachSubscriber(
		NewUpdateWorkSubscriber(a.dao),
		UpdateTopic)
	if err != nil {
		log.Errorf("Failed to attach subscriber to WorkEngine: %s", err.Error())
		os.Exit(1)
	}
	err = a.engine.AttachSubscriber(
		NewBindingWorkSubscriber(a.dao),
		BindingTopic)
	if err != nil {
		log.Errorf("Failed to attach subscriber to WorkEngine: %s", err.Error())
		os.Exit(1)
	}
	err = a.engine.AttachSubscriber(
		NewUnbindingWorkSubscriber(a.dao),
		UnbindingTopic)
	if err != nil {
		log.Errorf("Failed to attach subscriber to WorkEngine: %s", err.Error())
		os.Exit(1)
	}
	log.Debugf("Active work engine topics: %+v", a.engine.GetActiveTopics())

	apb.InitializeSecretsCache(a.config.GetSubConfigArray("secrets"))
	apb.InitializeClusterConfig(a.config.GetSubConfig("openshift"))

	rules := []rbac.PolicyRule{}
	if !a.brokerConfig.AutoEscalate {
		rules, err = retrieveClusterRoleRules(a.config.GetString("openshift.sandbox_role"))
		if err != nil {
			log.Errorf("Unable to retrieve cluster roles rules from cluster\n"+
				" You must be using OpenShift 3.7 to use the User rules check.\n%v", err)
			os.Exit(1)
		}
	}
	a.clusterRoleRules = rules

	return nil // TODO: Should be bubbling up errors so the framework can handle them gracefully
}

func validateRegistryNames(registrys []registries.Registry) {
	names := map[string]bool{}
	for _, registry := range registrys {
		if _, ok := names[registry.RegistryName()]; ok {
			panic(fmt.Sprintf("Name of registry: %v must be unique", registry.RegistryName()))
		}
		names[registry.RegistryName()] = true
	}
}

func retrieveClusterRoleRules(clusterRole string) ([]rbac.PolicyRule, error) {
	k8scli, err := clients.Kubernetes()
	if err != nil {
		return nil, err
	}

	// Retrieve Cluster Role that has been defined.
	k8sRole, err := k8scli.Client.RbacV1beta1().ClusterRoles().Get(clusterRole, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return convertAPIRbacToK8SRbac(k8sRole).Rules, nil
}

// convertAPIRbacToK8SRbac - because we are using the kubernetes validation,
// and they have not started using the authoritative api package for their own
// types, we need to do some conversion here now that we are on client-go 5.0.X
func convertAPIRbacToK8SRbac(apiRole *apirbac.ClusterRole) *rbac.ClusterRole {
	rules := []rbac.PolicyRule{}
	for _, pr := range apiRole.Rules {
		rules = append(rules, rbac.PolicyRule{
			Verbs:           pr.Verbs,
			APIGroups:       pr.APIGroups,
			Resources:       pr.Resources,
			ResourceNames:   pr.ResourceNames,
			NonResourceURLs: pr.NonResourceURLs,
		})
	}
	return &rbac.ClusterRole{
		TypeMeta:   apiRole.TypeMeta,
		ObjectMeta: apiRole.ObjectMeta,
		Rules:      rules,
	}
}

func initClients(c *config.Config) error {
	// Designed to panic early if we cannot construct required clients.
	// this likely means we're in an unrecoverable configuration or environment.
	// Best we can do is alert the operator as early as possible.
	//
	// Deliberately forcing the injection of deps here instead of running as a
	// method on the app. Forces developers at authorship time to think about
	// dependencies / make sure things are ready.
	log.Notice("Initializing clients...")
	log.Debug("Trying to connect to etcd")

	// Intialize the etcd configuration
	clients.InitEtcdConfig(c)
	etcdClient, err := clients.Etcd()
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	version, err := etcdClient.GetVersion(ctx)
	if err != nil {
		return err
	}

	log.Infof("Etcd Version [Server: %s, Cluster: %s]", version.Server, version.Cluster)

	_, err = clients.Kubernetes()
	if err != nil {
		return err
	}

	return nil
}
