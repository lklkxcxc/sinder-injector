package main

import (
	"flag"
	"os"

	hook "injector/webhook"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var log = logf.Log.WithName("example-controller")

type HookParamters struct {
	certDir       string
	port          int
}

func main() {
	var params HookParamters

	flag.IntVar(&params.port, "port", 8443, "Wehbook port")
	flag.StringVar(&params.certDir, "certDir", "/certs/", "Wehbook certificate folder")
	flag.Parse()

	logf.SetLogger(zap.Logger(false))
	entryLog := log.WithName("entrypoint")

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}


	// Setup webhooks
	entryLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	hookServer.Port = params.port
	hookServer.CertDir = params.certDir

	entryLog.Info("registering webhooks to the webhook server")
	hookServer.Register("/mutate", &webhook.Admission{Handler: &hook.SidecarInjector{Name: "Logger", Client: mgr.GetClient(), Annotations: map[string]string{"sidecar-injector": "true"}}})

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
