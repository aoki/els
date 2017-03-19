NAME			:= els
VERSION		:= v0.2.0
REVISION	:= $(shell git rev-parse --short HEAD)

SRCS		:= $(shell find src -type f -name '*.go')
LDFLAGS	:= -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

.DEFAULT_GOAL := bin/$(NAME)

.PHONY: dep-install
dep-install:
ifeq ($(shell command -v dep 2> /dev/null),)
	go get -u github.com/golang/dep/...
endif

.PHONY: init
init:
	dep init

.PHONY: deps
deps:
	dep ensure -update

bin/$(NAME): $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME) $(SRCS)

.PHONY: install
install:
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o ${GOPATH}/bin/$(NAME)

.PHONY: cross-build
cross-build: deps
	@for os in darwin linux windows; do \
		for arch in amd64 386; do \
			printf "Building for $${os} $${arch}\n";\
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME)  $(SRCS); \
		done; \
	done


.PHONEY: release
INPUT := {\"tag_name\": \"$(VERSION)\", \"target_commitish\": \"master\", \"draft\": false, \"prerelease\": false }
release: dist
	@printf "Generate artifacts\n"
	@mkdir -p artifact
	@\ls dist | xargs -IXXX zip -j artifact/$(NAME)_XXX.zip dist/XXX/$(NAME)
	@\ls dist | xargs -IXXX tar czf artifact/$(NAME)_XXX.tar.gz -C dist/XXX $(NAME)
	@printf "Create relese $(VERSION)\n"
	$(eval RELEASE_ID := $(shell curl --fail -s -X POST https://api.github.com/repos/$(GITHUB_USER)/$(NAME)/releases \
		-H "Accept: application/vnd.github.v3+json" \
		-H  "Authorization: token $(GITHUB_TOKEN)" \
		-H "Content-Type: application/json" \
		-d "$(INPUT)" | tr ',' '\n' | grep id | head -1 | cut -d':' -f2 | tr -d ' '))
	$(eval RELEASE_ID := 5789241)
	@for ARCHIVE in $$(\ls -ld1 $(PWD)/artifact/*); do \
		echo $${ARCHIVE}; \
		ARCHIVE_NAME=$$(basename $${ARCHIVE}); \
		CONTENT_TYPE=$$(file --mime-type -b $${ARCHIVE}); \
		curl --fail -X POST https://uploads.github.com/repos/$(GITHUB_USER)/$(NAME)/releases/$(RELEASE_ID)/assets?name=$${ARCHIVE_NAME} \
			-H "Accept: application/vnd.github.v3+json" \
			-H "Authorization: token $(GITHUB_TOKEN)" \
			-H "Content-Type: $${CONTENT_TYPE}" \
			--data-binary @"$${ARCHIVE}"; \
	done
