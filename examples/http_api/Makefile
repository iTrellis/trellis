.PHONY: build
build:
	tbuilder build --config=build.yaml

.PHONY: build_remote
build_remote:
	tbuilder build --config=build_remote.yaml

.PHONY: run
run: 
	./sample run --config=run.yaml

.PHONY: run_remote
run_remote: 
	./remote run --config=run_remote.yaml

.PHONY: clean
clean:
	rm sample remote