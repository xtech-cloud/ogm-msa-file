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
	MICRO_REGISTRY=consul micro call omo.api.msa.file Healthy.Echo '{"msg":"hello"}'
	# -------------------------------------------------------------------------
	# 创建存储桶, 缺少参数
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.Make
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.Make '{"name":"test1"}'
	# 创建存储桶，本地 ,10G (1024x1024x1024x10)
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.Make '{"name":"local", "capacity": 10737418240}'
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.Make '{"name":"minio", "capacity": 10737418240}'
	# 创建存储桶，已存在
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.Make '{"name":"local", "capacity": 10737418240}'
	# 列举存储桶，无参数
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.List
	# 列举存储桶
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.List '{"offset":1, "count":1}'
	# 更新存储桶引擎,无参数
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.UpdateEngine
	# 更新存储桶引擎,不存在
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.UpdateEngine '{"uuid":"000000000000000", "engine":2}'
	# 更新存储桶引擎
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.UpdateEngine '{"uuid":"132b6c5fc9193d6ae58027ae302ab67b", "engine":2, "address":"localhost:9000", "scope": "test", "accessKey": "root", "accessSecret":"minio@OMO"}'
	# 更新存储桶引擎
	# 更新存储桶容量,不存在
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.UpdateCapacity '{"uuid":"test"}'
	# 更新存储桶容量, 100G
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.UpdateCapacity '{"uuid":"132b6c5fc9193d6ae58027ae302ab67b", "capacity":107374182400}'
	# 更新token, 不存在
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.ResetToken '{"uuid":"test"}'
	# 更新token
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.ResetToken '{"uuid":"132b6c5fc9193d6ae58027ae302ab67b"}'
	# 获取存储桶，无参数
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.Get
	# 获取存储桶，不存在
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.Get '{"uuid":"00000000"}'
	# 获取存储桶
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.Get '{"uuid":"132b6c5fc9193d6ae58027ae302ab67b"}'
	# 准备对象元数据, 超过容量
	MICRO_REGISTRY=consul micro call omo.api.msa.file Object.Prepare '{"bucket":"132b6c5fc9193d6ae58027ae302ab67b", "uname":"cc2bd8f09bb88b5dd20f9b432631b8ca.jpg", "size":107374182401}'
	# 准备对象元数据
	MICRO_REGISTRY=consul micro call omo.api.msa.file Object.Prepare '{"bucket":"132b6c5fc9193d6ae58027ae302ab67b", "uname":"cc2bd8f09bb88b5dd20f9b432631b8ca.jpg", "size":223345}'
	# 写入对象元数据 
	MICRO_REGISTRY=consul micro call omo.api.msa.file Object.Flush '{"bucket":"132b6c5fc9193d6ae58027ae302ab67b", "uname":"cc2bd8f09bb88b5dd20f9b432631b8ca.jpg", "path":"a-1.jpg"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.file Object.Flush '{"bucket":"132b6c5fc9193d6ae58027ae302ab67b", "uname":"cc2bd8f09bb88b5dd20f9b432631b8ca.jpg", "path":"a-2.jpg"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.file Object.Flush '{"bucket":"132b6c5fc9193d6ae58027ae302ab67b", "uname":"cc2bd8f09bb88b5dd20f9b432631b8ca.jpg", "path":"a-3.jpg"}'
	# 列举对象 
	MICRO_REGISTRY=consul micro call omo.api.msa.file Object.List '{"bucket":"132b6c5fc9193d6ae58027ae302ab67b"}'
	# 列举对象 
	MICRO_REGISTRY=consul micro call omo.api.msa.file Object.List '{"bucket":"132b6c5fc9193d6ae58027ae302ab67b", "offset":1, "count":1}'
	# 获取存储桶
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.Get '{"uuid":"132b6c5fc9193d6ae58027ae302ab67b"}'
	# 删除存储桶，无参数
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.Remove
	# 删除存储桶，不存在
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.Remove '{"uuid":"00000"}'
	# 删除存储桶，存在
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.Remove '{"uuid":"132b6c5fc9193d6ae58027ae302ab67b"}'
	MICRO_REGISTRY=consul micro call omo.api.msa.file Bucket.Remove '{"uuid":"f5ddaf0ca7929578b408c909429f68f2"}'

.PHONY: tcall
tcall:
	mkdir -p ./bin
	go build -o ./bin/ ./tester
	./bin/tester

.PHONY: post
post:
	curl -X POST -d '{"msg":"hello"}' 127.0.0.1:8080/msa/file/Healthy/Echo

.PHONY: dist
dist:
	mkdir dist
	tar -zcf dist/${APP_NAME}-${BUILD_VERSION}.tar.gz ./bin/${APP_NAME}

.PHONY: docker
docker:
	docker build . -t omo-msa-startkit:latest
