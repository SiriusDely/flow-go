COLLECTION = 3
CONSENSUS = 3
VALID_CONSENSUS := $(shell test $(CONSENSUS) -ge 2; echo $$?)
EXECUTION = 2
VALID_EXECUTION := $(shell test $(EXECUTION) -ge 2; echo $$?)
VERIFICATION = 1
ACCESS = 2
OBSERVER = 0
NCLUSTERS=1
EPOCHLEN=10000   # 0 means use default
STAKINGLEN=2000 # 0 means use default
DKGLEN=2000     # 0 means use default
CONSENSUS_DELAY=800ms
COLLECTION_DELAY=950ms

PROFILER=false
PROFILE_UPLOADER=false
TRACING=true
EXTENSIVE_TRACING=false
CADENCE_TRACING=false

# The Git commit hash
COMMIT=$(shell git rev-parse HEAD)
VERSION=localnetbuild

.PHONY: init
init:
ifeq ($(strip $(VALID_EXECUTION)), 1)
	# multiple execution nodes are required to prevent seals being generated in case of execution forking.
	$(error Number of Execution nodes should be no less than 2)
else ifeq ($(strip $(VALID_CONSENSUS)), 1)
	$(error Number of Consensus nodes should be no less than 2)
else
	go run -tags relic \
		-ldflags="-X 'github.com/onflow/flow-go/cmd/build.commit=${COMMIT}' \
		-X  'github.com/onflow/flow-go/cmd/build.semver=${VERSION}'" \
		bootstrap.go \
		-collection=$(COLLECTION) \
		-consensus=$(CONSENSUS) \
		-execution=$(EXECUTION) \
		-verification=$(VERIFICATION) \
		-access=$(ACCESS) \
		-observer=$(OBSERVER) \
		-nclusters=$(NCLUSTERS) \
		-epoch-length=$(EPOCHLEN) \
		-epoch-staking-phase-length=$(STAKINGLEN) \
		-epoch-dkg-phase-length=$(DKGLEN) \
		-profiler=$(PROFILER) \
		-profile-uploader=$(PROFILE_UPLOADER) \
		-tracing=$(TRACING) \
		-cadence-tracing=$(CADENCE_TRACING) \
		-extensive-tracing=$(EXTENSIVE_TRACING) \
		-consensus-delay=$(CONSENSUS_DELAY) \
		-collection-delay=$(COLLECTION_DELAY)
endif

# Creates a light version of the localnet with just 1 instance for each node type
.PHONY: init-light
init-light:
	$(MAKE) -e COLLECTION=1 CONSENSUS=2 EXECUTION=1 VERIFICATION=1 ACCESS=1 NCLUSTERS=1 init

# Creates a version of localnet configured with short epochs
.PHONY: init-short-epochs
init-short-epochs:
	$(MAKE) -e EPOCHLEN=200 STAKINGLEN=10 DKGLEN=50 init

.PHONY: start
start:
	DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 docker-compose -f docker-compose.metrics.yml up -d --remove-orphans
	DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 docker-compose -f docker-compose.nodes.yml build
	DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 docker-compose -f docker-compose.nodes.yml up -d

.PHONY: stop
stop:
	DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 docker-compose -f docker-compose.metrics.yml -f docker-compose.nodes.yml down -v --remove-orphans

.PHONY: load
load:
	go run --tags relic ../benchmark/cmd/manual -log-level info -tps 1,10,100 -tps-durations 30s,30s

.PHONY: tps-test
tps-test:
	go run --tags relic ../benchmark/cmd/manual -log-level info -tps 25 -tps-durations 1m

.PHONY: tps-ci-smoke
tps-ci-smoke:
	go run --tags relic ../benchmark/cmd/ci -log-level info -tps 10 -tps-durations 2m -slice-size 10s

.PHONY: tps-test-ci
tps-test-ci:
	make init
	make start
	go run --tags relic ../benchmark/cmd/ci -tps 1
	make stop
