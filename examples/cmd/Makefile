.PHONY: build
build:
	tbuilder build --config=build.yaml

.PHONY: build_remote
build_remote:
	tbuilder build --config=build_remote.yaml

.PHONY: build_grpc
build_grpc:
	tbuilder build --config=build_grpc.yaml

.PHONY: run
run: 
	./sample run --config=run.yaml

.PHONY: run_remote
run_remote: 
	./remote run --config=run_remote.yaml

.PHONY: run_remote2
run_remote2: 
	./sample run --config=run_remote2.yaml

.PHONY: run_grpc
run_grpc: 
	./grpc run --config=run_grpc.yaml

.PHONY: clean
clean:
	rm sample remote
