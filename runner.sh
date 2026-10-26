#!/bin/bash

REBUILD=${REBUILD:=0}
APP=mimic3-api
VERSION=0.1.0
IMAGE=mimic3-api:0.1.0
DOCKERFILE=Dockerfile

ROOT_PATH=$(git rev-parse --show-toplevel)
if [ -n "$ROOT_PATH" ]; then
	cd "$ROOT_PATH"
fi

siginit_handler(){
	printf "> Terminating:\n"
	docker ps | grep $APP > /dev/null 2>&1
	if [ $? -eq 0 ]; then
		printf "> Killing app:\n"
		docker exec -ti $APP sh -c "ps -A | grep $APP | awk '{ print \$1 }' | xargs kill -s 2"
		sleep 2
	fi
	docker ps | grep $APP > /dev/null 2>&1
	if [ $? -eq 0 ]; then
		printf "> Stoping container\n"
		docker stop $APP
		if [ $? -ne 0 ]; then
			printf "> ERROR: Failed to stop container\n"
			exit -1
		fi
	fi
	printf "> Done\n"
}

# trap siginit_handler SIGINT
if [ $? -ne 0 ]; then
	printf "> ERROR: Failed to trap SIGINT:\n"
	exit 1
fi

rm -f bin/$APP

if [ "$REBUILD" == "1" ]; then
	docker rmi $APP:$VERSION > /dev/null 2>/dev/null
fi

printf "> Building new image:\n"
docker images --format "{{.Repository}}:{{.Tag}}" | grep "$APP:$VERSION" > /dev/null 2>&1
if [ $? -ne 0 ]; then
	docker build \
		--platform linux/amd64 \
		--progress plain \
		--build-arg APP=$APP \
		--tag "$APP:$VERSION" \
		--file $DOCKERFILE \
		.

	if [ $? -ne 0 ]; then
		printf "> ERROR: Build new image:\n"
		exit 1
	fi
fi

printf "> Starting container:\n"
docker run -ti --rm \
	--platform linux/amd64 \
	--volume $PWD/tmp/data:/data \
	--name $APP \
	$APP:$VERSION

