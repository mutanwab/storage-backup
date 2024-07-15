package main

import (
	storagev1 "github.com/mutanwab/storage-backup/pkg/apis/storage.backup.io/v1beta1"
	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
	"os"
)

func main() {
	os.Unsetenv("GOPATH")
	controllergen.Run(args.Options{
		OutputPackage: "storage-backup/pkg/generated",
		Boilerplate:   "scripts/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"storage.backup.io": {
				Types: []interface{}{
					storagev1.StorageBackup{},
				},
				GenerateTypes:   true,
				GenerateClients: true,
			},
		},
	})
}
