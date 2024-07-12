package main

import (
	"os"

	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
	v1 "k8s.io/api/core/v1"
	"storage-backup/pkg/codegen/generator"
)

func main() {
	os.Unsetenv("GOPATH")

	controllergen.Run(args.Options{
		OutputPackage: "github.com/rancher/rancher/pkg/generated",
		Boilerplate:   "scripts/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"cluster.x-k8s.io": {
				Types: []interface{}{},
			},
		},
	})
	generator.GenerateNativeTypes(v1.SchemeGroupVersion, []interface{}{
		v1.Pod{},
	}, nil)
}
