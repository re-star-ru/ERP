# include .env
# .DEFAULT_GOAL := dev

dev:
	hugo server -D --bind="0.0.0.0"

build:
	hugo

docker: docker image build .