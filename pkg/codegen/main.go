package main

import (
	"os"
	"storage-backup/pkg/apis/storage.backup.io/v1beta1"

	"codegen/generator"
	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
	v1 "k8s.io/api/core/v1"
)

func main() {
	os.Unsetenv("GOPATH")

	controllergen.Run(args.Options{
		OutputPackage: "storage-backup/pkg/generated",
		Boilerplate:   "scripts/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"storage.backup.io": {
				PackageName: "storage.backup.io",
				Types: []interface{}{
					// All structs with an embedded ObjectMeta field will be picked up
					"./pkg/apis/storage.backup.io/v1beta1",
					v1beta1.StorageBackup{},
				},
				GenerateTypes: true,
			},
		},
	})

	generator.GenerateNativeTypes(v1.SchemeGroupVersion, []interface{}{
		v1.Pod{},
	}, nil)
}
