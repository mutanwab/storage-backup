package backup

import (
	"fmt"
	bTypes "github.com/longhorn/backupstore/types"
	butil "github.com/longhorn/backupstore/util"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

type ProgressState string

const (
	snapBlockSize = 2 << 20 // 2MiB

	ProgressStateInProgress = ProgressState("in_progress")
	ProgressStateComplete   = ProgressState("complete")
	ProgressStateError      = ProgressState("error")
)

type VolumeSnapStatus struct {
	lock            sync.Mutex
	Name            string
	volumeID        string
	SnapshotID      string
	SnapFileHandler *os.File
	Error           string
	Progress        int
	BackupURL       string
	State           ProgressState
	IsIncremental   bool
	IsOpened        bool
}

func NewVolumeSnap(backupName, volumeName, snapshotName string) *VolumeSnapStatus {
	if backupName == "" {
		backupName = butil.GenerateName("backup")
	}
	return &VolumeSnapStatus{
		Name:       backupName,
		State:      ProgressStateInProgress,
		volumeID:   volumeName,
		SnapshotID: snapshotName,
	}
}

func (v *VolumeSnapStatus) HasSnapshot(snapID, volumeID string) bool {
	v.lock.Lock()
	defer v.lock.Unlock()
	if v.volumeID != volumeID {
		logrus.Warnf("Invalid state volume [%s] are open, not [%s]", v.volumeID, volumeID)
		return false
	}
	snapPath := GetSnapshotPath(snapID)

	if _, err := os.Stat(snapPath); os.IsNotExist(err) {
		logrus.Errorf("snap file %s is not exist", snapPath)
		return false
	}

	return true
}

func (v *VolumeSnapStatus) CompareSnapshot(snapID, compareID, volumeID string) (*bTypes.Mappings, error) {
	v.lock.Lock()
	if err := v.assertOpen(snapID, volumeID); err != nil {
		v.lock.Unlock()
		return nil, err
	}
	v.lock.Unlock()

	mappings := &bTypes.Mappings{
		BlockSize: snapBlockSize,
	}
	mapping := bTypes.Mapping{
		Offset: -1,
	}

	snapPath := GetSnapshotPath(snapID)
	fileSize, err := GetAvailableSpaceBlock(snapPath)
	if err != nil {
		return nil, err
	}
	count := fileSize / snapBlockSize
	for i := int64(0); i < count; i++ {
		offset := i * snapBlockSize
		if offset > fileSize {
			offset = fileSize
		}
		mapping = bTypes.Mapping{
			Offset: offset,
			Size:   snapBlockSize,
		}
		mappings.Mappings = append(mappings.Mappings, mapping)
	}

	return mappings, nil
}

func (v *VolumeSnapStatus) OpenSnapshot(snapID, volumeID string) error {
	var err error
	id := GetSnapshotPath(snapID)
	logrus.Infof("GenerateSnapshotDiskName get snap disk name: %s", id)
	v.SnapFileHandler, err = os.Open(id)
	if err != nil {
		return err
	}
	v.IsOpened = true
	return nil
}

func (v *VolumeSnapStatus) ReadSnapshot(snapID, volumeID string, start int64, data []byte) error {
	v.lock.Lock()
	defer v.lock.Unlock()
	if err := v.assertOpen(snapID, volumeID); err != nil {
		return err
	}

	_, err := v.SnapFileHandler.ReadAt(data, start)
	return err
}

func (v *VolumeSnapStatus) CloseSnapshot(snapID, volumeID string) error {
	err := v.SnapFileHandler.Close()
	if err != nil {
		return err
	}
	v.IsOpened = false
	return nil
}

func (v *VolumeSnapStatus) UpdateBackupStatus(snapID, volumeID string, backupState string, backupProgress int, backupURL string, err string) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if !v.isVolumeSnapshotMatched(snapID, volumeID) {
		return fmt.Errorf("invalid volume [%s] and snapshot [%s], not volume [%s], snapshot [%s]", v.volumeID, v.SnapshotID, volumeID, snapID)
	}

	v.State = ProgressState(backupState)
	v.Progress = backupProgress
	v.BackupURL = backupURL
	v.Error = err

	if v.Progress == 100 {
		v.State = ProgressStateComplete
	} else if v.Error != "" {
		v.State = ProgressStateError
	}
	return nil
}

func (v *VolumeSnapStatus) isVolumeSnapshotMatched(id, volumeID string) bool {
	if v.volumeID != volumeID || v.SnapshotID != id {
		return false
	}
	return true
}

func (v *VolumeSnapStatus) assertOpen(id, volumeID string) error {
	if v.volumeID != volumeID || v.SnapshotID != id {
		return fmt.Errorf("invalid volume [%s] and snapshot [%s], not volume [%s], snapshot [%s]", v.volumeID, v.SnapshotID, volumeID, id)
	}
	if !v.IsOpened {
		return fmt.Errorf("volume [%s] and snapshot [%s] are not opened", volumeID, id)
	}
	return nil
}
