build:
	go build

build-linux:
	GOOS=linux GOARCH=386 go build

build-docker:
	docker build -t csdread/arcam-controller:2 .

push:
	docker push csdread/arcam-controller:2

clean:
	rm -rf arcam-controller
