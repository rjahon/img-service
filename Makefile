CURRENT_DIR=$(shell pwd)

gen-proto-module:
	./script/genproto.sh ${CURRENT_DIR}

rm-proto-module:
	sudo rm -rf genproto

migration-up:
	migrate -path migrations -database 'postgres://postgres:postgres@0.0.0.0:5432/img_service?sslmode=disable' up

migration-down:
	migrate -path migrations -database 'postgres://postgres:postgres@0.0.0.0:5432/img_service?sslmode=disable' down

run:
	go run cmd/main.go