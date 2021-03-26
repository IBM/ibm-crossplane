#!/bin/bash
#
# Copyright 2021 IBM Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

############################################################
# GKE section
############################################################
PROJECT ?= oceanic-guard-191815
ZONE    ?= us-west1-a
CLUSTER ?= prow

activate-serviceaccount:
ifdef GOOGLE_APPLICATION_CREDENTIALS
	gcloud auth activate-service-account --key-file="$(GOOGLE_APPLICATION_CREDENTIALS)"
endif

get-cluster-credentials: activate-serviceaccount
	gcloud container clusters get-credentials "$(CLUSTER)" --project="$(PROJECT)" --zone="$(ZONE)"

config-docker: get-cluster-credentials
	@ibm/scripts/config_docker.sh

############################################################
# Prow section
############################################################

# Specify whether this repo is build locally or not, default values is '1';
# If set to 1, then you need to also set 'DOCKER_USERNAME' and 'DOCKER_PASSWORD'
# environment variables before build the repo.
BUILD_LOCALLY ?= 1
IMAGE_NAME ?= ibm-crossplane
RELEASE_VERSION ?= $(shell cat RELEASE_VERSION)
GO_SUPPORTED_VERSIONS = 1.14|1.15

export OSBASEIMAGE=registry.access.redhat.com/ubi8/ubi-minimal:latest

ifeq ($(BUILD_LOCALLY),0)
	DOCKER_REGISTRY = hyc-cloud-private-integration-docker-local.artifactory.swg-devops.com/ibmcom
	export BUILD_REGISTRY=$(DOCKER_REGISTRY)
else
	DOCKER_REGISTRY = hyc-cloud-private-scratch-docker-local.artifactory.swg-devops.com/ibmcom
endif

ifeq ($(HOSTOS),darwin)
	MANIFEST_TOOL_ARGS ?= --username $(DOCKER_USERNAME) --password $(DOCKER_PASSWORD)
else
	MANIFEST_TOOL_ARGS ?= 
endif

image-amd64:
	docker tag $(BUILD_REGISTRY)/crossplane-amd64:latest $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION)-amd64
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION)-amd64

image-ppc64le:
	docker tag $(BUILD_REGISTRY)/crossplane-ppc64le:latest $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION)-ppc64le
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION)-ppc64le

image-s390x:
	docker tag $(BUILD_REGISTRY)/crossplane-s390x:latest $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION)-s390x
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION)-s390x

images: $(MANIFEST_TOOL)
ifeq ($(BUILD_LOCALLY),1)
	@make build
	@make image-amd64
	@$(MANIFEST_TOOL) $(MANIFEST_TOOL_ARGS) push from-args --platforms linux/amd64 --template $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION)-ARCH --target $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION) || $(FAIL)
else
	@make config-docker
	@make build.all
	@make image-amd64
	@make image-ppc64le
	@make image-s390x
	@$(MANIFEST_TOOL) $(MANIFEST_TOOL_ARGS) push from-args --platforms linux/amd64,linux/ppc64le,linux/s390x --template $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION)-ARCH --target $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION) || $(FAIL)
endif


