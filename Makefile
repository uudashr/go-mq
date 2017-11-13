SOURCES := $(shell find . -name '*.go' -type f -not -path './vendor/*'  -not -path '*/mocks/*')

NSQ_NSQD_ADDR ?= 127.0.0.1:4150
NSQ_LOOKUPD_ADDR ?= 127.0.0.1:4161
NSQ_TOPIC ?= greet
NSQ_CHANNEL ?= public

# Dependencies Management
.PHONY: vendor-prepare
vendor-prepare:
	@echo "Installing glide"
	@curl https://glide.sh/get | sh

glide.lock: glide.yaml
	@glide update

.PHONY: vendor-update
vendor-update:
	@glide update

vendor: glide.lock
	@glide install

.PHONY: clean-vendor
clean-vendor:
	@rm -rf vendor

# Testing
.PHONY: test
test: vendor
	@go test -short $$(glide novendor)

.PHONY: test-nsq
test-nsq: vendor
	@go test -v ./nsq -nsq.nsqd-addr "$(NSQ_NSQD_ADDR)" -nsq.lookupd-addr "$(NSQ_LOOKUPD_ADDR)" -nsq.topic "$(NSQ_TOPIC)" -nsq.channel "$(NSQ_CHANNEL)"

.PHONY: test-pubsub
test-pubsub: vendor
	@go test -v ./pubsub -gcp.project-id "$(GCP_PROJECT)" -gcp.topic-id "$(GCP_TOPIC)" -gcp.subscription-id "$(GCP_SUBSCRIPTION)" -gcp.credentials-file "$(GCP_CREDENTIALS_FILE)"

# Upstream service
.PHONY: docker-nsq-up
docker-nsq-up:
	@docker-compose -f docker-compose-nsq.yml up -d && docker-compose -f docker-compose-nsq.yml logs -f

.PHONY: docker-nsq-down
docker-nsq-down:
	@docker-compose -f docker-compose-nsq.yml down
