

build-data-image:
	docker build -t test-data-image -f Dockerfile.data .

run-data:
	docker run -it --name test-data-container test-data-image


