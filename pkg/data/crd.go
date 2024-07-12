package data

import (
	"context"
	"k8s.io/client-go/rest"

	//loggingv1 "github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"

	"storage-backup/pkg/util/crd"
)

func createCRDs(ctx context.Context, restConfig *rest.Config) error {
	factory, err := crd.NewFactoryFromClient(ctx, restConfig)
	if err != nil {
		return err
	}
	return factory.
		BatchCreateCRDsIfNotExisted(
		//crd.NonNamespacedFromGV(loggingv1.GroupVersion, "Logging", loggingv1.Logging{}),
		).
		BatchCreateCRDsIfNotExisted(
		//crd.FromGV(lhv1beta2.SchemeGroupVersion, "Backup", lhv1beta2.Backup{}),
		).
		BatchWait()
}
