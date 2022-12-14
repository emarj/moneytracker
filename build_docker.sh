#!/bin/bash

VERSION=$(git describe --tags --abbrev=8 --dirty --always --long)
REPO="emarj/moneytracker:"
TAG="$REPO$VERSION"
LATEST="${REPO}latest"
BUILD_TIMESTAMP=$( date '+%F_%H:%M:%S' )
docker build --platform linux/amd64 -t "$TAG" -t "$LATEST" --build-arg VERSION="$VERSION" --build-arg BUILD_TIMESTAMP="$BUILD_TIMESTAMP" . 
docker push "$TAG" 
docker push "$LATEST"