REGISTRY?=registry.fly.io

.PHONY: deploy
deploy: bin/ko bin/flyctl
	$(eval img=$(shell KO_DOCKER_REPO=${REGISTRY} bin/ko publish --sbom none --base-import-paths ./cmd/spartan))
	bin/flyctl deploy --image ${img}

include Makefile.bindl
