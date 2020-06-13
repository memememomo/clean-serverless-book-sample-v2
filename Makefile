DOCKER_YAML=-f docker-compose.yml
DOCKER=COMPOSE_PROJECT_NAME=corporate-blog-recommend docker-compose $(DOCKER_YAML)

build:
	$(DOCKER) build ${ARGS}

docker-up:
	$(DOCKER) up

go-lint:
	$(DOCKER) run go-test ./scripts/go-lint.sh

go-test:
	$(DOCKER) run go-test ./scripts/go-test.sh '${PACKAGE}' '${ARGS}'

go-build:
	$(DOCKER) run go-test ./scripts/build-handlers.sh

go-get:
	$(DOCKER) run go-test go get ${ARGS}

sls-package:
	$(DOCKER) run sls sls package

npm-install:
	$(DOCKER) run sls npm install ${ARGS}

deploy: go-build
	$(DOCKER) run sls ./scripts/deploy.sh

delete-stack:
	$(DOCKER) run sls sls remove
