SHELL=bash

ENV:=dev
# Hard variables
DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
OUR_BASE_FILE=${DIR}/deployment/base.yaml
OUR_ENV_FILE=${DIR}/deployment/${ENV}.yaml
TPL_FILE=${DIR}/deployment/job.nomadtpl
# Vary variables
BASE_PATH:=.
BASE_FILE:=${BASE_PATH}/deployment/base.yaml
ENV_VAR_FILE:=${BASE_PATH}/deployment/${ENV}.yaml
OUT_FILE:=

.PHONY: job.nomad
job.nomad:
ifeq ("${OUT_FILE}", "")
	$(eval OUT_FILE := /tmp/job.$(shell uuidgen).nomad)
endif
	$(eval var_file := $(shell mktemp /tmp/env.$(shell uuidgen).yaml))
	@yq merge -x ${OUR_BASE_FILE} ${OUR_ENV_FILE} ${BASE_FILE} ${ENV_VAR_FILE} > ${var_file}
	@cd ${BASE_PATH} && levant render -var-file=${var_file} -out=${OUT_FILE} ${TPL_FILE} 2>/dev/null
	@echo "${OUT_FILE}"

