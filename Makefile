BINARY_NAME_LOMS=loms-service


build-loms:
	go build -o $(BINARY_NAME_LOMS) ./cmd/server

run-loms:
	./$(BINARY_NAME_LOMS)


run-all: build-loms run-loms


# Проверка покрытия
check-coverage:
	go test -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | findstr total
	del coverage.out


# Когнитивная нагрузка
cognitive-load:
	gocognit -top 10 -ignore "_mock|_test" .\internal


# Цикломатическая нагрузка
cyclomatic-load:
	gocyclo -top 10 -ignore "_mock|_test" .\internal


# Используем bin в текущей директории для установки protoc
LOCAL_BIN := $(CURDIR)/bin

# Добавляем bin в текущей директории в РАТН при запуске protoc
PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc


# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
	https://github.com/protocolbuffers/protobuf vendor-proto/protobuf && \
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

# Удаление папки vendor-proto
.PHONY: .vendor-rm
.vendor-rm:
	rm -rf vendor-proto


# Установка бинарный зависимостей
.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies ... )
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
	mv $(LOCAL_BIN)/protoc-gen-go.exe $(LOCAL_BIN)/protoc-gen-go && \
	mv $(LOCAL_BIN)/protoc-gen-go-grpc.exe $(LOCAL_BIN)/protoc-gen-go-grpc

.vendor-proto: vendor-proto/google/protobuf


NOTES_PROTO_PATH := "api/loms/v1"

.make-dir:
	mkdir -p pkg/api/loms/v1

.PHONY: .protoc-generate
.protoc-generate: .bin-deps .vendor-proto .make-dir
	protoc \
	-I ${NOTES_PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out pkg/${NOTES_PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out pkg/${NOTES_PROTO_PATH} \
    --go-grpc_opt paths=source_relative \
	api/loms/v1/loms.proto
	go mod tidy
