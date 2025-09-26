.PHONY: start build wire update new clean
PROJECT_NAME := synapse
SERVICE_NAME := $(PROJECT_NAME)
MODULE_NAME := github.com/xh-polaris/$(PROJECT_NAME)

HANDLER_DIR := biz/api/controller
MODEL_DIR := biz/api/model
ROUTER_DIR := biz/api/router

IDL_DIR ?= ../$(PROJECT_NAME)-idl/
FULL_MAIN_IDL_PATH := $(IDL_DIR)/api.thrift

IDL_OPTIONS := -I $(IDL_DIR) --idl $(FULL_MAIN_IDL_PATH)
OUTPUT_OPTIONS := --handler_dir $(HANDLER_DIR) --model_dir $(MODEL_DIR) --router_dir $(ROUTER_DIR)
EXTRA_OPTIONS := --pb_camel_json_tag=true --unset_omitempty=true --enable_extends=true

run:
	sh ./output/bootstrap.sh
build:
	sh ./build.sh
build_and_run:
	sh ./build.sh && sh ./output/bootstrap.sh
wire:
	wire ./provider
update:
	hz --verbose update $(IDL_OPTIONS) --mod $(MODULE_NAME) $(EXTRA_OPTIONS)
	@files=$$(find biz/application/dto -type f); \
	for file in $$files; do \
  	  sed -i  -e 's/func init\(\).*//' $$file; \
  	done
new:
	hz new $(IDL_OPTIONS) $(OUTPUT_OPTIONS) --service $(SERVICE_NAME) --module $(MODULE_NAME) $(EXTRA_OPTIONS)
clean:
	rm -r ./output
wire:
	wire gen ./provider