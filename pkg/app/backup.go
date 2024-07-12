package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"storage-backup/pkg/backup"
)

func BackupCmd() cli.Command {
	return cli.Command{
		Name: "backup",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "backup-name",
				Required: true,
				Usage:    "Specify backup name",
			},
			cli.StringFlag{
				Name:     "volume-name",
				Required: true,
				Usage:    "Specify volume name",
			},
			cli.StringFlag{
				Name:     "snap-name",
				Required: true,
				Usage:    "Specify snap name",
			},
			cli.StringFlag{
				Name:     "backup-target",
				Required: true,
				Usage:    "Specify backup target url",
			},
		},
		Action: func(c *cli.Context) {
			if err := startBackup(c); err != nil {
				logrus.Fatalf("Error starting manager: %v", err)
			}
		},
	}
}

func startBackup(c *cli.Context) error {
	backupName := c.String("backup-name")
	volumeName := c.String("volume-name") //pvc-id
	snapshotFileName := c.String("snap-name")
	target := c.String("backup-target")
	backupTarget := fmt.Sprintf("%s?backup=%s&volume=%s", target, backupName, volumeName)
	var labels []string
	backupStatus, backupConfig, err := backup.DoBackupInit(backupName, volumeName, snapshotFileName, backupTarget, labels)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to initialize backup %v", backupName)
		return err
	}

	err = backup.DoBackupCreate(backupStatus, backupConfig)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to create backup %v", backupName)
		return err
	}
	return nil
}
