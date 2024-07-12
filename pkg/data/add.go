package data

import (
	"context"

	"storage-backup/pkg/config"
)

// Init adds built-in resources
func Init(ctx context.Context, mgmtCtx *config.Management, options config.Options) error {
	if err := createCRDs(ctx, mgmtCtx.RestConfig); err != nil {
		return err
	}

	if err := addPublicNamespace(mgmtCtx.Apply); err != nil {
		return err
	}
	// Not applying the built-in templates and secrets in case users have edited them.
	//if err := createTemplates(mgmtCtx, publicNamespace); err != nil {
	//	return err
	//}
	return nil
}
