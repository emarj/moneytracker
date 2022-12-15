#!/bin/sh
BRANCH=$(git branch --show-current)
VERSION=$(git describe --tags --dirty --always)
REPO="emarj/moneytracker"

if [ $BRANCH == "master" ];
then
    BRANCH=""
else
    BRANCH="${BRANCH}-"
fi

TAGGED="${REPO}:${BRANCH}${VERSION}"
LATEST="${REPO}:${BRANCH}latest"
BUILD_TIMESTAMP=$( date '+%F_%H:%M:%S' )

docker build --platform linux/amd64 -t "$TAGGED" -t "$LATEST" --build-arg VERSION="$VERSION" --build-arg BUILD_TIMESTAMP="$BUILD_TIMESTAMP" . 
