module codegen

go 1.20

replace (
	github.com/docker/distribution => github.com/docker/distribution v2.7.1+incompatible // oras dep requires a replace is set
	github.com/docker/docker => github.com/docker/docker v20.10.9+incompatible // oras dep requires a replace is set

	github.com/knative/pkg => github.com/rancher/pkg v0.0.0-20190514055449-b30ab9de040e
	github.com/matryer/moq => github.com/rancher/moq v0.0.0-20200712062324-13d1f37d2d77

	github.com/opencontainers/runc => github.com/opencontainers/runc v1.1.2
	github.com/rancher/rancher/pkg/apis => ./pkg/apis
	github.com/rancher/rancher/pkg/client => ./pkg/client
	github.com/rancher/steve => github.com/mutanwab/steve v1.1.0

	github.com/tencentyun/tcecloud-sdk-go v0.0.0-incompatible => ./pkg/gientech/clusterhandler/tcecloud-sdk-go

	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc => go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.20.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp => go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.20.0
	go.opentelemetry.io/otel => go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/otlp => go.opentelemetry.io/otel/exporters/otlp v0.20.0
	go.opentelemetry.io/otel/sdk => go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/otel/trace => go.opentelemetry.io/otel/trace v0.20.0
	go.opentelemetry.io/proto/otlp => go.opentelemetry.io/proto/otlp v0.7.0

	helm.sh/helm/v3 => github.com/rancher/helm/v3 v3.9.0-rancher1
	k8s.io/api => k8s.io/api v0.24.5
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.24.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.24.5
	k8s.io/apiserver => k8s.io/apiserver v0.24.5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.24.5
	k8s.io/client-go => github.com/rancher/client-go v1.24.0-rancher1
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.24.5
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.24.5
	k8s.io/code-generator => k8s.io/code-generator v0.24.5
	k8s.io/component-base => k8s.io/component-base v0.24.5
	k8s.io/component-helpers => k8s.io/component-helpers v0.24.5
	k8s.io/controller-manager => k8s.io/controller-manager v0.24.5
	k8s.io/cri-api => k8s.io/cri-api v0.24.5
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.24.5
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.24.5
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.24.5
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.24.5
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.24.5
	k8s.io/kubectl => k8s.io/kubectl v0.24.5
	k8s.io/kubelet => k8s.io/kubelet v0.24.5
	k8s.io/kubernetes => k8s.io/kubernetes v1.24.2
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.24.5
	k8s.io/metrics => k8s.io/metrics v0.24.5
	k8s.io/mount-utils => k8s.io/mount-utils v0.24.5
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.24.5
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.24.5

	sigs.k8s.io/aws-iam-authenticator => github.com/rancher/aws-iam-authenticator v0.5.9-0.20220713170329-78acb8c83863
	sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v1.2.0
)

require (
	github.com/rancher/norman v0.0.0-20220627222520-b74009fac3ff
	github.com/rancher/wrangler v0.6.2-0.20200820173016-2068de651106
	k8s.io/api v0.25.6
	k8s.io/apimachinery v0.24.5
	k8s.io/gengo v0.0.0-20211129171323-c02415ce4185
)

require (
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/matryer/moq v0.0.0-20200607124540-4638a53893e6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.10-0.20220218145154-897bd77cd717 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/code-generator v0.24.5 // indirect
	k8s.io/klog/v2 v2.60.1 // indirect
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9 // indirect
	sigs.k8s.io/json v0.0.0-20211208200746-9f7c6b3444d2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
)
