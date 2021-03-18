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
# Prow section
############################################################

# Specify whether this repo is build locally or not, default values is '1';
# If set to 1, then you need to also set 'DOCKER_USERNAME' and 'DOCKER_PASSWORD'
# environment variables before build the repo.
BUILD_LOCALLY ?= 1
IMAGE_NAME ?= ibm-crossplane
RELEASE_VERSION ?= $(shell cat RELEASE_VERSION)

ifeq ($(BUILD_LOCALLY),0)
    export CONFIG_DOCKER_TARGET = config-docker
	DOCKER_REGISTRY = hyc-cloud-private-integration-docker-local.artifactory.swg-devops.com/ibmcom
else
	DOCKER_REGISTRY = hyc-cloud-private-scratch-docker-local.artifactory.swg-devops.com/ibmcom
endif

images: $(MANIFEST_TOOL)
	docker tag $(BUILD_REGISTRY)/crossplane-amd64:latest $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION)-amd64
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION)-amd64
ifeq ($(HOSTOS),darwin)
	@$(MANIFEST_TOOL) --username $(DOCKER_USERNAME) --password $(DOCKER_PASSWORD) push from-args --platforms linux/amd64 --template $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION)-ARCH --target $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION) || $(FAIL)
else
	@$(MANIFEST_TOOL) push from-args --platforms linux/amd64 --template $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION)-ARCH --target $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(RELEASE_VERSION) || $(FAIL)
endif

############################################################
# GKE section
############################################################
PROJECT ?= oceanic-guard-191815
ZONE    ?= us-west1-a
CLUSTER ?= prow
GLCOUD ?= $(shell which gcloud)

activate-serviceaccount:
	$(GCLOUD) auth activate-service-account --key-file="$(GOOGLE_APPLICATION_CREDENTIALS)"

get-cluster-credentials: activate-serviceaccount
	$(GCLOUD) container clusters get-credentials "$(CLUSTER)" --project="$(PROJECT)" --zone="$(ZONE)"

config-docker: get-cluster-credentials
	@common/scripts/config_docker.sh

