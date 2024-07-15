package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mutanwab/storage-backup/pkg/config"

	"github.com/rancher/dynamiclistener"
	"github.com/rancher/dynamiclistener/server"
	"github.com/rancher/wrangler/pkg/ratelimit"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type BackupServer struct {
	Context       context.Context
	RESTConfig    *restclient.Config
	DynamicClient dynamic.Interface
	ClientSet     *kubernetes.Clientset

	Handler http.Handler
}

func New(ctx context.Context, clientConfig clientcmd.ClientConfig, options config.Options) (*BackupServer, error) {
	var err error
	BackupServer := &BackupServer{
		Context: ctx,
	}
	BackupServer.RESTConfig, err = clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	BackupServer.RESTConfig.RateLimiter = ratelimit.None

	if err := Wait(ctx, BackupServer.RESTConfig); err != nil {
		return nil, err
	}

	BackupServer.ClientSet, err = kubernetes.NewForConfig(BackupServer.RESTConfig)
	if err != nil {
		return nil, fmt.Errorf("kubernetes clientset create error: %s", err.Error())
	}

	BackupServer.DynamicClient, err = dynamic.NewForConfig(BackupServer.RESTConfig)
	if err != nil {
		return nil, fmt.Errorf("kubernetes dynamic client create error:%s", err.Error())
	}

	return BackupServer, nil
}

func Wait(ctx context.Context, config *rest.Config) error {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	for {
		_, err := client.Discovery().ServerVersion()
		if err == nil {
			break
		}
		logrus.Infof("Waiting for server to become available: %v", err)
		select {
		case <-ctx.Done():
			return fmt.Errorf("startup canceled")
		case <-time.After(2 * time.Second):
		}
	}

	return nil
}

func (s *BackupServer) ListenAndServe(listenerCfg *dynamiclistener.Config, opts config.Options) error {
	listenOpts := &server.ListenOpts{}

	if listenerCfg != nil {
		listenOpts.TLSListenerConfig = *listenerCfg
	}

	if err := server.ListenAndServe(s.Context, opts.HTTPSListenPort, opts.HTTPListenPort, s.Handler, listenOpts); err != nil {
		return err
	}

	<-s.Context.Done()
	return s.Context.Err()
}
