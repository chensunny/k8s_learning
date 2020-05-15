#!/bin/sh
set -x -e

buildNumber=""
if [[ -f "build.json" ]];then
buildNumber="`cat build.json`"
fi

nameFlag="-X github.com/AfterShip/golang-common/whoami.name=CRMProcessor" \
&& versionFlag="-X github.com/AfterShip/golang-common/whoami.version="`cat version`"" \
&& buildNumberFlag="-X github.com/AfterShip/golang-common/whoami.buildNumber="${buildNumber}"" \
&& buildAtFlag="-X github.com/AfterShip/golang-common/whoami.buildAt="`date '+%Y-%m-%dT%H:%M:%S%z'`"" \
&& buildOnFlag="-X github.com/AfterShip/golang-common/whoami.buildOn="`hostname`"" \
&& gitCommitFlag="-X github.com/AfterShip/golang-common/whoami.gitCommit="`git rev-parse HEAD`@`git name-rev --name-only HEAD`"" \
&& ldflags="${nameFlag} ${versionFlag} ${buildNumberFlag} ${buildAtFlag} ${buildOnFlag} ${gitCommitFlag}" \
&& go build  -ldflags "${ldflags}" -mod=vendor ./cmd/web

#can build manual in pod
#go build -mod=vendor -o ./cmd/organization/organization ./cmd/organization
#go build -mod=vendor -o ./cmd/processor/processor ./cmd/processor
