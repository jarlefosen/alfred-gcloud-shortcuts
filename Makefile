
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
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-s -w' -trimpath -o bin/projects cmd/projects/*.go

build-products:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-s -w' -trimpath -o bin/products cmd/products/*.go

build: build-projects build-products

sort-products:
	cat products.json | jq -s '.[] | sort_by(.name)' > products_sorted.json
	cp products_sorted.json products.json
	rm products_sorted.json
