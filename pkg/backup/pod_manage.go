package backup

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"text/template"
)

var (
	PodYaml = template.Must(template.New("cluster").Parse(BackupPodTemplate))
)

type BackupPod struct {
	Pvcs    []string `json:"pvcs"`
	Ns      string   `json:"ns"`
	PodType string   `json:"podType"`
}

// 启动pod挂载pvc, 备份pvc

// 删除pod

func GenerateYaml(pod *BackupPod) (string, error) {
	// 生成pod模板
	buf := &bytes.Buffer{}
	if err := PodYaml.Execute(buf, pod); err != nil {
		logrus.Errorf("installCluster: generate sealer kubefile error: %s", err.Error())
		return "", err
	}
	return buf.String(), nil
}
