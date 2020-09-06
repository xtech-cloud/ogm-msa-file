APP_NAME := omo-msa-file
BUILD_VERSION   := $(shell git tag --contains)
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD )

.PHONY: build
build:
	go build -ldflags \
		"\
		-X 'main.BuildVersion=${BUILD_VERSION}' \
		-X 'main.BuildTime=${BUILD_TIME}' \
		-X 'main.CommitID=${COMMIT_SHA1}' \
		"\
		-o ./bin/${APP_NAME}

.PHONY: run
run:
	./bin/${APP_NAME}

.PHONY: install
install:
	go install

.PHONY: clean
clean:
	rm -rf /tmp/msa-file.db

.PHONY: call
call:
	# -------------------------------------------------------------------------
	# 创建存储桶, 缺少参数
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Make
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Make '{"name":"test1"}'
	# 创建存储桶，本地 ,10G (1024x1024x1024x10)
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Make '{"name":"local", "capacity": 10737418240}'
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Make '{"name":"qiniu", "capacity": 10737418240}'
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Make '{"name":"minio", "capacity": 10737418240}'
	# 创建存储桶，已存在
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Make '{"name":"local", "capacity": 10737418240}'
	# 列举存储桶，无参数
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.List
	# 列举存储桶
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.List '{"offset":1, "count":1}'
	# 更新存储桶引擎,无参数
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.UpdateEngine
	# 更新存储桶引擎,不存在
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.UpdateEngine '{"name":"test", "engine":2}'
	# 更新存储桶引擎
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.UpdateEngine '{"name":"minio", "engine":2, "address":"localhost:9000", "scope": "test", "accessKey": "root", "accessSecret":"minio@OMO"}'
	# 更新存储桶引擎
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.UpdateEngine '{"name":"qiniu", "engine":3, "address":"ugc.meex.tech", "scope": "meex-ugc", "accessKey": "OM-zigYG0s0HH2od-KBIoRSRo2L90ZPuZ0vqKQPj", "accessSecret":"nup6a_1QXRD4oktAtKs_9lrUF44r652WT5IBxrrH"}'
	# 更新存储桶容量,不存在
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.UpdateCapacity '{"name":"test"}'
	# 更新存储桶容量, 100G
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.UpdateCapacity '{"name":"minio", "capacity":107374182400}'
	# 更新token, 不存在
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.ResetToken '{"name":"test"}'
	# 更新token
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.ResetToken '{"name":"minio"}'
	# 获取存储桶，无参数
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Get
	# 获取存储桶，不存在
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Get '{"name":"test"}'
	# 获取存储桶
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Get '{"name":"qiniu"}'
	# 获取凭证
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Auth '{"name":"qiniu"}'
	# 删除存储桶，无参数
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Remove
	# 删除存储桶，不存在
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Remove '{"name":"test"}'
	# 删除存储桶，存在
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Remove '{"name":"minio"}'
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Remove '{"name":"qiniu"}'
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Remove '{"name":"local"}'

.PHONY: tcall
tcall:
	mkdir -p ./bin
	go build -o ./bin/ ./tester
	./bin/tester

.PHONY: dist
dist:
	mkdir dist
	tar -zcf dist/${APP_NAME}-${BUILD_VERSION}.tar.gz ./bin/${APP_NAME}

.PHONY: docker
docker:
	docker build . -t omo-msa-startkit:latest
