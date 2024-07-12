package config

import (
	"context"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/start"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type (
	_scaledKey struct{}
)

type Options struct {
	Namespace       string
	Threadiness     int
	HTTPListenPort  int
	HTTPSListenPort int

	RancherEmbedded bool
	RancherURL      string
	HCIMode         bool
}

type Scaled struct {
	Ctx               context.Context
	ControllerFactory controller.SharedControllerFactory
	starters          []start.Starter

	Management *Management
}

type Management struct {
	ctx               context.Context
	Apply             apply.Apply
	ControllerFactory controller.SharedControllerFactory

	ClientSet  *kubernetes.Clientset
	RestConfig *rest.Config

	starters []start.Starter
}

func SetupScaled(ctx context.Context, restConfig *rest.Config) (context.Context, *Scaled, error) {
	scaled := &Scaled{
		Ctx: ctx,
	}
	var err error
	scaled.Management, err = setupManagement(ctx, restConfig)
	if err != nil {
		return nil, nil, err
	}

	return context.WithValue(scaled.Ctx, _scaledKey{}, scaled), scaled, nil
}

func setupManagement(ctx context.Context, restConfig *rest.Config) (*Management, error) {
	management := &Management{
		ctx: ctx,
	}

	apply, err := apply.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	management.Apply = apply

	return management, nil
}

func ScaledWithContext(ctx context.Context) *Scaled {
	return ctx.Value(_scaledKey{}).(*Scaled)
}

func (s *Scaled) Start(threadiness int) error {
	return start.All(s.Ctx, threadiness, s.starters...)
}

func (s *Management) Start(threadiness int) error {
	return start.All(s.ctx, threadiness, s.starters...)
}
