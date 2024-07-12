package server

import (
	"github.com/rancher/wrangler/pkg/schemes"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var (
	localSchemeBuilder = runtime.SchemeBuilder{}
	AddToScheme        = localSchemeBuilder.AddToScheme
	Scheme             = runtime.NewScheme()
)

func init() {
	utilruntime.Must(AddToScheme(Scheme))
	utilruntime.Must(schemes.AddToScheme(Scheme))
}
