#!/bin/bash

# updates commonInitialisms from most recent go lint source.

SRC=$(realpath $(cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd ))

PKG="github.com/golang/lint"

OUT=$SRC/initialisms.go

set -e

# get latest revision
if [ "$1" == "--update" ]; then
	go get -u $PKG
fi

# get git rev info
PKGPATH="$(go env GOPATH)/src/$PKG"
pushd $PKGPATH &> /dev/null
PKGVER=$(git rev-parse --short master)
popd &> /dev/null

DATA="package $(basename $SRC)

// commonInitialisms is the set of commonInitialisms.
//
// taken from: $PKG @ $PKGVER
$(sed '/{/{:1; /}/!{N; b1}; /var commonInitialisms/p}; d' $PKGPATH/lint.go)"

echo "$DATA" > $OUT

gofmt -w $OUT
