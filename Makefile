build:
	go build

build-linux:
	GOOS=linux GOARCH=386 go build

build-docker:
	docker build -t csdread/arcam-controller:1 .

push:
	docker push csdread/arcam-controller:1

clean:
	rm -rf arcam-controller
