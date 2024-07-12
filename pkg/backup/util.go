package backup

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	MaximumVolumeNameSize = 64

	blockdevFileName = "/usr/sbin/blockdev"
	SnapDir          = ""

	VolumeHeadDiskPrefix = "volume-head-"
	VolumeHeadDiskSuffix = ".img"
	VolumeHeadDiskName   = VolumeHeadDiskPrefix + "%03d" + VolumeHeadDiskSuffix

	SnapshotDiskPrefix = "volume-snap-"
	SnapshotDiskSuffix = ".img"
	SnapshotDiskName   = SnapshotDiskPrefix + "%s" + SnapshotDiskSuffix

	DeltaDiskPrefix = "volume-delta-"
	DeltaDiskSuffix = ".img"
	DeltaDiskName   = DeltaDiskPrefix + "%s" + DeltaDiskSuffix

	DiskMetadataSuffix = ".meta"
	DiskChecksumSuffix = ".checksum"

	snapTmpSuffix = ".snap_tmp"

	expansionSnapshotInfix = "expand-%d"

	replicaExpansionLabelKey = "replica-expansion"
)

var (
	validVolumeName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]+$`)
)

// GetAvailableSpaceBlock gets the amount of available space at the block device path specified.
func GetAvailableSpaceBlock(deviceName string) (int64, error) {
	// Check if the file exists and is a device file.
	if ok, err := IsDevice(deviceName); !ok || err != nil {
		return int64(-1), err
	}

	// Device exists, attempt to get size.
	cmd := exec.Command(blockdevFileName, "--getsize64", deviceName)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return int64(-1), fmt.Errorf("%v, %s", err, errBuf.String())
	}
	i, err := strconv.ParseInt(strings.TrimSpace(out.String()), 10, 64)
	if err != nil {
		return int64(-1), err
	}
	return i, nil
}

func GetFileSize(fileName string) (int64, error) {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return int64(-1), err
	}
	return fileInfo.Size(), nil
}

// IsDevice returns true if it's a device file
func IsDevice(deviceName string) (bool, error) {
	info, err := os.Stat(deviceName)
	if err == nil {
		return (info.Mode() & os.ModeDevice) != 0, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func GenerateSnapshotDiskName(name string) string {
	return fmt.Sprintf(SnapshotDiskName, name)
}

func GetSnapshotPath(name string) string {
	return filepath.Join("/dev", name)
}

func GenerateSnapshotDiskChecksumName(diskName string) string {
	return diskName + DiskChecksumSuffix
}

func GenerateSnapshotDiskMetaName(diskName string) string {
	return diskName + DiskMetadataSuffix
}

func GenerateDeltaFileName(name string) string {
	return fmt.Sprintf(DeltaDiskName, name)
}

func GenerateSnapTempFileName(fileName string) string {
	return fileName + snapTmpSuffix
}

func GetSnapshotNameFromTempFileName(tmpFileName string) (string, error) {
	if !strings.HasSuffix(tmpFileName, snapTmpSuffix) {
		return "", fmt.Errorf("invalid snapshot tmp filename")
	}
	return strings.TrimSuffix(tmpFileName, snapTmpSuffix), nil
}

func GetSnapshotNameFromDiskName(diskName string) (string, error) {
	if !strings.HasPrefix(diskName, SnapshotDiskPrefix) || !strings.HasSuffix(diskName, SnapshotDiskSuffix) {
		return "", fmt.Errorf("invalid snapshot disk name %v", diskName)
	}
	result := strings.TrimPrefix(diskName, SnapshotDiskPrefix)
	result = strings.TrimSuffix(result, SnapshotDiskSuffix)
	return result, nil
}

func GenerateExpansionSnapshotName(size int64) string {
	return fmt.Sprintf(expansionSnapshotInfix, size)
}

func GenerateExpansionSnapshotLabels(size int64) map[string]string {
	return map[string]string{
		replicaExpansionLabelKey: strconv.FormatInt(size, 10),
	}
}

func IsHeadDisk(diskName string) bool {
	if strings.HasPrefix(diskName, VolumeHeadDiskPrefix) &&
		strings.HasSuffix(diskName, VolumeHeadDiskSuffix) {
		return true
	}
	return false
}

func ValidVolumeName(name string) bool {
	if len(name) > MaximumVolumeNameSize {
		return false
	}
	return validVolumeName.MatchString(name)
}

func ParseLabels(labels []string) (map[string]string, error) {
	result := map[string]string{}
	for _, label := range labels {
		kv := strings.SplitN(label, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid label not in <key>=<value> format %v", label)
		}
		key := kv[0]
		value := kv[1]
		if errList := IsQualifiedName(key); len(errList) > 0 {
			return nil, fmt.Errorf("invalid key %v for label: %v", key, errList[0])
		}
		// We don't need to validate the Label value since we're allowing for any form of data to be stored, similar
		// to Kubernetes Annotations. Of course, we should make sure it isn't empty.
		if value == "" {
			return nil, fmt.Errorf("invalid empty value for label with key %v", key)
		}
		result[key] = value
	}
	return result, nil
}

const qnameCharFmt string = "[A-Za-z0-9]"
const qnameExtCharFmt string = "[-A-Za-z0-9_.]"
const qualifiedNameFmt = "(" + qnameCharFmt + qnameExtCharFmt + "*)?" + qnameCharFmt
const qualifiedNameErrMsg string = "must consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character"
const qualifiedNameMaxLength int = 63

var qualifiedNameRegexp = regexp.MustCompile("^" + qualifiedNameFmt + "$")

// IsQualifiedName tests whether the value passed is what Kubernetes calls a
// "qualified name".  This is a format used in various places throughout the
// system.  If the value is not valid, a list of error strings is returned.
// Otherwise an empty list (or nil) is returned.
func IsQualifiedName(value string) []string {
	var errs []string
	parts := strings.Split(value, "/")
	var name string
	switch len(parts) {
	case 1:
		name = parts[0]
	case 2:
		var prefix string
		prefix, name = parts[0], parts[1]
		if len(prefix) == 0 {
			errs = append(errs, "prefix part "+EmptyError())
		} else if msgs := IsDNS1123Subdomain(prefix); len(msgs) != 0 {
			errs = append(errs, prefixEach(msgs, "prefix part ")...)
		}
	default:
		return append(errs, "a qualified name "+RegexError(qualifiedNameErrMsg, qualifiedNameFmt, "MyName", "my.name", "123-abc")+
			" with an optional DNS subdomain prefix and '/' (e.g. 'example.com/MyName')")
	}

	if len(name) == 0 {
		errs = append(errs, "name part "+EmptyError())
	} else if len(name) > qualifiedNameMaxLength {
		errs = append(errs, "name part "+MaxLenError(qualifiedNameMaxLength))
	}
	if !qualifiedNameRegexp.MatchString(name) {
		errs = append(errs, "name part "+RegexError(qualifiedNameErrMsg, qualifiedNameFmt, "MyName", "my.name", "123-abc"))
	}
	return errs
}

const dns1123LabelFmt string = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"

var dns1123LabelRegexp = regexp.MustCompile("^" + dns1123LabelFmt + "$")

const dns1123SubdomainFmt = dns1123LabelFmt + "(\\." + dns1123LabelFmt + ")*"
const dns1123SubdomainErrorMsg string = "a DNS-1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character"

// DNS1123SubdomainMaxLength is a subdomain's max length in DNS (RFC 1123)
const DNS1123SubdomainMaxLength int = 253

var dns1123SubdomainRegexp = regexp.MustCompile("^" + dns1123SubdomainFmt + "$")

// IsDNS1123Subdomain tests for a string that conforms to the definition of a
// subdomain in DNS (RFC 1123).
func IsDNS1123Subdomain(value string) []string {
	var errs []string
	if len(value) > DNS1123SubdomainMaxLength {
		errs = append(errs, MaxLenError(DNS1123SubdomainMaxLength))
	}
	if !dns1123SubdomainRegexp.MatchString(value) {
		errs = append(errs, RegexError(dns1123SubdomainErrorMsg, dns1123SubdomainFmt, "example.com"))
	}
	return errs
}

// MaxLenError returns a string explanation of a "string too long" validation
// failure.
func MaxLenError(length int) string {
	return fmt.Sprintf("must be no more than %d characters", length)
}

// RegexError returns a string explanation of a regex validation failure.
func RegexError(msg string, fmt string, examples ...string) string {
	if len(examples) == 0 {
		return msg + " (regex used for validation is '" + fmt + "')"
	}
	msg += " (e.g. "
	for i := range examples {
		if i > 0 {
			msg += " or "
		}
		msg += "'" + examples[i] + "', "
	}
	msg += "regex used for validation is '" + fmt + "')"
	return msg
}

// EmptyError returns a string explanation of a "must not be empty" validation
// failure.
func EmptyError() string {
	return "must be non-empty"
}

func prefixEach(msgs []string, prefix string) []string {
	for i := range msgs {
		msgs[i] = prefix + msgs[i]
	}
	return msgs
}
