package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/longhorn/backupstore"
	"github.com/longhorn/backupstore/types"
	"github.com/longhorn/backupstore/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"runtime"
	"runtime/debug"
	"time"
)

var (
	VERSION = "0.0.0"
	log     = logrus.WithFields(logrus.Fields{"pkg": "backup"})
)

type ErrorResponse struct {
	Error string
}

func ResponseLogAndError(v interface{}) {
	if e, ok := v.(*logrus.Entry); ok {
		e.Error(e.Message)
		fmt.Println(e.Message)
	} else {
		e, isErr := v.(error)
		_, isRuntimeErr := e.(runtime.Error)
		if isErr && !isRuntimeErr {
			logrus.Errorf(fmt.Sprint(e))
			fmt.Println(fmt.Sprint(e))
		} else {
			logrus.Errorf("Caught FATAL error: %s", v)
			debug.PrintStack()
			fmt.Println("Caught FATAL error: ", v)
		}
	}
}

// ResponseOutput would generate a JSON format byte array of object for output
func ResponseOutput(v interface{}) ([]byte, error) {
	j, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return j, nil
}

func RequiredMissingError(name string) error {
	return fmt.Errorf("cannot find valid required parameter: %v", name)
}

func DoBackupInit(backupName, volumeName, snapshotName, destURL string, labels []string) (*VolumeSnapStatus, *backupstore.DeltaBackupConfig, error) {
	log.Infof("Initializing backup %v for volume %v snapshot %v", backupName, volumeName, snapshotName)

	var (
		err      error
		labelMap map[string]string
	)

	if volumeName == "" || snapshotName == "" || destURL == "" {
		return nil, nil, fmt.Errorf("missing input parameter")
	}

	if !ValidVolumeName(volumeName) {
		return nil, nil, fmt.Errorf("invalid volume name %v for backup %v", volumeName, backupName)
	}

	if labels != nil {
		labelMap, err = ParseLabels(labels)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "cannot parse backup labels for backup %v", backupName)
		}
	}

	snapPath := GetSnapshotPath(snapshotName)
	size, err := GetAvailableSpaceBlock(snapPath)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "get backup file %s size error", snapPath)
	} else if size == int64(-1) {
		size, err = GetFileSize(snapPath)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "get backup file %s size error", snapPath)
		} else if size == int64(-1) {
			return nil, nil, errors.Wrapf(err, "get backup file %s size error", snapPath)
		}
	}

	backup := NewVolumeSnap(backupName, volumeName, snapshotName)

	volume := &backupstore.Volume{
		Name:        volumeName,
		Size:        size,
		Labels:      labelMap,
		CreatedTime: util.Now(),
	}
	snapshot := &backupstore.Snapshot{
		Name:        snapshotName,
		CreatedTime: util.Now(),
	}

	log.Debugf("Starting backup for %v, snapshot %v, dest %v", volume, snapshot, destURL)
	config := &backupstore.DeltaBackupConfig{
		BackupName:      backupName,
		Volume:          volume,
		Snapshot:        snapshot,
		DestURL:         destURL,
		DeltaOps:        backup,
		Labels:          labelMap,
		ConcurrentLimit: 1,
	}

	return backup, config, nil
}

func DoBackupCreate(volumeSnapStatus *VolumeSnapStatus, config *backupstore.DeltaBackupConfig) error {
	log.Infof("Start creating backup %v", volumeSnapStatus.Name)

	isIncremental, err := backupstore.CreateDeltaBlockBackup(config.BackupName, config)
	if err != nil {
		return err
	}
	for {
		if string(volumeSnapStatus.State) == string(types.ProgressStateInProgress) {
			time.Sleep(5 * time.Second)
		} else {
			time.Sleep(5 * time.Second)
			break
		}
	}

	volumeSnapStatus.IsIncremental = isIncremental
	logrus.Infof("do backup end")
	return nil
}

func DoBackupRestore(backupURL string, toFile string, restoreObj *RestoreStatus) error {
	backupURL = util.UnescapeURL(backupURL)
	log.Debugf("Start restoring from %v into snapshot %v", backupURL, toFile)

	config := &backupstore.DeltaRestoreConfig{
		BackupURL:       backupURL,
		DeltaOps:        restoreObj,
		Filename:        toFile,
		ConcurrentLimit: int32(1),
	}
	ctx := context.Background()

	if err := backupstore.RestoreDeltaBlockBackup(ctx, config); err != nil {
		return err
	}
	for {
		if string(restoreObj.State) == string(types.ProgressStateInProgress) {
			time.Sleep(5 * time.Second)
		} else {
			time.Sleep(5 * time.Second)
			break
		}
	}
	logrus.Infof("restore %s end", backupURL)

	return nil
}

func DoBackupRestoreIncrementally(url string, deltaFile string, lastRestored string,
	restoreObj *RestoreStatus) error {
	backupURL := util.UnescapeURL(url)
	log.Debugf("Start incremental restoring from %v into delta file %v", backupURL, deltaFile)

	config := &backupstore.DeltaRestoreConfig{
		BackupURL:      backupURL,
		DeltaOps:       restoreObj,
		LastBackupName: lastRestored,
		Filename:       deltaFile,
	}
	ctx := context.Background()

	if err := backupstore.RestoreDeltaBlockBackupIncrementally(ctx, config); err != nil {
		return err
	}

	return nil
}
