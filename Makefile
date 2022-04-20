build:
	go build

build-linux:
	GOOS=linux GOARCH=386 go build

build-docker:
	docker build -t csdread/arcam-controller:1 .

push:
	docker push csdread/arcam-controller:1

deploy:
	kubectl -n default apply -f k8s/deployment.yaml
	kubectl -n default apply -f k8s/service.yaml

clean:
	kubectl -n default delete -f k8s/deployment.yaml
	rm -rf arcam-controller
