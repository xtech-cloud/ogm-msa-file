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
	# 创建存储桶
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Make
	# 创建存储桶，本地
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Make '{"name":"test1"}'
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Make '{"name":"test2"}'
	# 创建存储桶，外部
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Make '{"name":"test3", "domain":"localhost:9000", "accessKey": "admin", "accessSecret":"password"}'
	# 创建存储桶，已存在
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Make '{"name":"test1"}'
	# 列举存储桶，无参数
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.List 
	# 列举存储桶
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.List '{"offset":1, "count":1}'
	# 更新存储桶,无参数
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Update
	# 更新存储桶
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Update '{"name":"test3", "domain":"localhost:9000", "accessKey": "root", "accessSecret":"minio@OMO"}'
	# 获取存储桶，无参数
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Get
	# 获取存储桶，不存在
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Get '{"name":"test"}'
	# 获取存储桶
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Get '{"name":"test3"}'
	# 删除存储桶，无参数
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Remove
	# 删除存储桶，不存在
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Remove '{"name":"test4"}'
	# 删除存储桶，存在
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Remove '{"name":"test3"}'
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Remove '{"name":"test2"}'
	MICRO_REGISTRY=consul micro call omo.msa.file Bucket.Remove '{"name":"test1"}'

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
