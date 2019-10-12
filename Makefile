
TARGET_DIR := target
WORKFLOW_FILE := $(TARGET_DIR)/alfred-gcloud-shortcuts.alfredworkflow


target/:
	mkdir target

clean:
	@[ -d $(TARGET_DIR) ] && rm -r $(TARGET_DIR) || true

workflow: build clean target/
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
