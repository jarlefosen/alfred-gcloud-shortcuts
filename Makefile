
BUILD_DIR := build
WORKFLOW_FILE := $(BUILD_DIR)/alfred-gcloud-shortcuts.alfredworkflow


build/:
	mkdir build

clean:
	@[ -d $(BUILD_DIR) ] && rm -r $(BUILD_DIR) || true

workflow: clean build/
	zip $(WORKFLOW_FILE) \
	info.plist \
	icon.png \
	products.json \
	bin/products \
	bin/projects

build-projects:
	GOOS=darwin GOARCH=amd64 go build -ldflags='-s -w' -o bin/projects cmd/projects/*.go

build-products:
	GOOS=darwin GOARCH=amd64 go build -ldflags='-s -w' -o bin/products cmd/products/*.go

build: build-projects build-products
