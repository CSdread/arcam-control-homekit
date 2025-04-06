build:
	go build

build-linux:
	GOOS=linux GOARCH=386 go build

build-docker: check-tag
	docker build --platform linux/amd64 -t csdread/arcam-controller:$(TAG) .

push: check-tag
	docker push csdread/arcam-controller:$(TAG)

ci: build-linux build-docker push

clean:
	rm -rf arcam-controller

check-tag:
ifndef TAG
	$(error TAG is undefined)
endif
