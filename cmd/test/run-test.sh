#!/bin/sh

docker-compose \
	-p sap-segmentation-test \
	up \
	--abort-on-container-exit \
	--exit-code-from sap-segmentation-test \
	--remove-orphans \
	--build \
	--force-recreate
docker-compose \
	-p sap-segmentation-test \
	down \
	--rmi local \
	--remove-orphans

