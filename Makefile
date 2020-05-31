default: build

build:
	go build .

install:
	go install .

clean:
	rm GoCloudTTS go.sum || true
