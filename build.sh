#!/bin/sh
DATE=$(date)
BRANCH=$(git branch --show-current)
COMMIT=$(git describe --tags --dirty --always --long)

LDFLAGS=
LDFLAGS="$LDFLAGS -X 'main.Commit=$COMMIT'"
LDFLAGS="$LDFLAGS -X 'main.Date=$DATE'"
LDFLAGS="$LDFLAGS -X 'main.Branch=$BRANCH'"

OSS=$(go env GOOS | xargs)
ARCH=$(go env GOARCH)

if [ $1 != "run" ];
then
    DEST="-o $2_${OSS}_${ARCH}"
fi

go $1 -v -ldflags "$LDFLAGS" $DEST cmd/server/main.go
