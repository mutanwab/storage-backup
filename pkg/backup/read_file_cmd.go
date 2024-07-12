package backup

//
//import (
//	"github.com/pkg/errors"
//	"github.com/sirupsen/logrus"
//	"io"
//	"os"
//)
//
//func main() {
//	snapName := "snap-file"
//	blkFileName := "/dev/cdi-block-volume"
//	if err := ReadDataFromBlkToFile(snapName, blkFileName); err != nil {
//		logrus.Errorf("read file %s err: %s", blkFileName, err.Error())
//		return
//	}
//}
//
//// ReadDataFromBlkToFile 从pvc读取数据到文件
//func ReadDataFromBlkToFile(wFileName, blkFileName string) error {
//	f, err := os.Open(wFileName)
//	if err != nil {
//		logrus.Errorf("open file %s err: %s", wFileName, err.Error())
//		return err
//	}
//	defer f.Close()
//	outFile, err := OpenFileOrBlockDevice(blkFileName)
//	if err != nil {
//		return err
//	}
//	defer outFile.Close()
//	logrus.Info("Reading data from pvc...\n")
//	if _, err = io.Copy(f, outFile); err != nil {
//		logrus.Errorf("Unable to write file from block file: %v\n", err)
//		return errors.Wrapf(err, "unable to write to file")
//	}
//	err = f.Sync()
//	return err
//}
//
//// StreamDataToFile 从文件读取数据到pvc
//// StreamDataToFile provides a function to stream the specified io.Reader to the specified local file
//func StreamDataToFile(r io.Reader, blkFileName string) error {
//	outFile, err := OpenFileOrBlockDevice(blkFileName)
//	if err != nil {
//		return err
//	}
//	defer outFile.Close()
//	logrus.Info("Writing data to block...\n")
//	if _, err = io.Copy(outFile, r); err != nil {
//		logrus.Errorf("Unable to write file from dataReader: %v\n", err)
//		os.Remove(outFile.Name())
//		return errors.Wrapf(err, "unable to write to file")
//	}
//	err = outFile.Sync()
//	return err
//}
//
//// OpenFileOrBlockDevice opens the destination data file, whether it is a block device or regular file
//func OpenFileOrBlockDevice(fileName string) (*os.File, error) {
//	var outFile *os.File
//	blockSize, err := GetAvailableSpaceBlock(fileName)
//	if err != nil {
//		return nil, errors.Wrapf(err, "error determining if block device exists")
//	}
//	if blockSize >= 0 {
//		// Block device found and size determined.
//		outFile, err = os.OpenFile(fileName, os.O_RDWR, os.ModePerm)
//	} else {
//		// Attempt to create the file with name filePath.  If it exists, fail.
//		outFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_EXCL|os.O_RDWR, os.ModePerm)
//	}
//	if err != nil {
//		return nil, errors.Wrapf(err, "could not open file %q", fileName)
//	}
//	return outFile, nil
//}
//
////func GenerateMetaFile(filePath string) error {
////
////}
