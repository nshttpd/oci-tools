
VERSION=`cat VERSION`
SHORTSHA=`git rev-parse --short HEAD`

LDFLAGS=-X github.com/nshttpd/oci-tools/cmd/oci-tool/cmd.version=$(VERSION)
LDFLAGS+=-X github.com/nshttpd/oci-tools/cmd/oci-tool/cmd.shortSha=$(SHORTSHA)

build:
	go build -ldflags "$(LDFLAGS)" .

utils:
	go get github.com/mitchellh/gox
	go get github.com/tcnksm/ghr

deploy: utils
	gox -parallel=4 -ldflags "$(LDFLAGS)" -output "dist/oci-tool_{{.OS}}_{{.Arch}}"
	ghr -t $(GITHUB_TOKEN) -u $(CIRCLE_PROJECT_USERNAME) -r $(CIRCLE_PROJECT_REPONAME) --replace $(VERSION) dist/