ROOT_DIR := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))/../

BALENA_BASE_TAG ?= 1.15-3.14-run
ORGANISATION ?= dirkmoors
BASE_IMAGE_NAME ?= nexus

.PHONY: build-raspberrypi3-64
build-raspberrypi3-64:
	docker ${DOCKER_ARGS} build -f ${ROOT_DIR}Dockerfile.raspberrypi3x64 -t ${ORGANISATION}/${BASE_IMAGE_NAME}:raspberrypi3-64-${BALENA_BASE_TAG} ${ROOT_DIR}

.PHONY: push-raspberrypi3-64
push-raspberrypi3-64: build-raspberrypi3-64
	docker ${DOCKER_ARGS} push ${ORGANISATION}/${BASE_IMAGE_NAME}:raspberrypi3-64-${BALENA_BASE_TAG}

.PHONY: build-raspberrypi4x64
build-raspberrypi4x64:
	docker ${DOCKER_ARGS} build -f ${ROOT_DIR}Dockerfile.raspberrypi4x64 -t ${ORGANISATION}/${BASE_IMAGE_NAME}:raspberrypi4-64-${BALENA_BASE_TAG} ${ROOT_DIR}

.PHONY: push-raspberrypi4x64
push-raspberrypi4x64: build-raspberrypi4x64
	docker ${DOCKER_ARGS} push ${ORGANISATION}/${BASE_IMAGE_NAME}:raspberrypi4-64-${BALENA_BASE_TAG}

.PHONY: build-raspberrypicm4-ioboard
build-raspberrypicm4-ioboard:
	docker ${DOCKER_ARGS} build -f ${ROOT_DIR}Dockerfile.raspberrypicm4-ioboard -t ${ORGANISATION}/${BASE_IMAGE_NAME}:raspberrypicm4-ioboard-${BALENA_BASE_TAG} ${ROOT_DIR}

.PHONY: push-raspberrypicm4-ioboard
push-raspberrypicm4-ioboard: build-raspberrypicm4-ioboard
	docker ${DOCKER_ARGS} push ${ORGANISATION}/${BASE_IMAGE_NAME}:raspberrypicm4-ioboard-${BALENA_BASE_TAG}

.PHONY: build-intel-nuc
build-intel-nuc:
	docker ${DOCKER_ARGS} build --platform linux/amd64 -f ${ROOT_DIR}Dockerfile.intel-nuc -t ${ORGANISATION}/${BASE_IMAGE_NAME}:intel-nuc-${BALENA_BASE_TAG} ${ROOT_DIR}

.PHONY: push-intel-nuc
push-intel-nuc: build-intel-nuc
	docker ${DOCKER_ARGS} push ${ORGANISATION}/${BASE_IMAGE_NAME}:intel-nuc-${BALENA_BASE_TAG}

.PHONY: pull-base-images
pull-base-images:
	docker ${DOCKER_ARGS} pull balenalib/generic-alpine-golang:1.15-3.14-run
	docker ${DOCKER_ARGS} pull balenalib/raspberrypi3-64-alpine-golang:1.15-3.14-run
	docker ${DOCKER_ARGS} pull balenalib/raspberrypi4-64-alpine-golang:1.15-3.14-run
	docker ${DOCKER_ARGS} pull balenalib/raspberrypi4-64-alpine-golang:1.15-3.14-run

.PHONY: build
build: build-intel-nuc build-raspberrypi4x64 build-raspberrypi3-64 build-raspberrypicm4-ioboard

.PHONY: push
push: push-intel-nuc push-raspberrypi4x64 push-raspberrypi3-64 push-raspberrypicm4-ioboard
