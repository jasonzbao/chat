build:
	go1.18 build -tags jsoniter -o cmd_dyna cmd/main.go

build_docker: 
	go build -tags jsoniter -o cmd_dyna cmd/main.go

r_entrypoint:
	./cmd_dyna --config ${CONFIG_FILE}
