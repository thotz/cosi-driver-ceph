module github.com/thotz/cosi-driver-ceph

go 1.15

require (
	github.com/QuentinPerez/go-radosgw v0.0.0-20210421143448-ef2be6928c41
	github.com/aws/aws-sdk-go v1.38.24
	github.com/google/uuid v1.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	google.golang.org/grpc v1.37.0
	k8s.io/apimachinery v0.19.4
	k8s.io/klog/v2 v2.8.0
	sigs.k8s.io/container-object-storage-interface-provisioner-sidecar v0.0.0-20210415211500-cb8b1286bb3c
	sigs.k8s.io/container-object-storage-interface-spec v0.0.0-20210330184956-b0de747ccee4
)
