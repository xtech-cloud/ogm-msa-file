module ogm-msa-file

go 1.16

require (
	github.com/asim/go-micro/plugins/config/encoder/yaml/v3 v3.0.0-20210721080634-e1bc7e302871
	github.com/asim/go-micro/plugins/config/source/etcd/v3 v3.0.0-20210721080634-e1bc7e302871
	github.com/asim/go-micro/plugins/logger/logrus/v3 v3.0.0-20210721080634-e1bc7e302871
	github.com/asim/go-micro/plugins/registry/etcd/v3 v3.0.0-20210721080634-e1bc7e302871
	github.com/asim/go-micro/plugins/server/grpc/v3 v3.0.0-20210721080634-e1bc7e302871
	github.com/asim/go-micro/v3 v3.5.2
	github.com/containerd/containerd v1.5.4 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/minio/minio-go/v7 v7.0.12
	github.com/opencontainers/selinux v1.8.2 // indirect
	github.com/qiniu/api.v7/v7 v7.8.2
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/xtech-cloud/ogm-msp-file v3.0.0+incompatible
	google.golang.org/genproto v0.0.0-20210721163202-f1cecdd8b78a // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/driver/mysql v1.1.1
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.12
)
