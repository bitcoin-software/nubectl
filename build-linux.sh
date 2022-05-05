#!/bin/sh
export PATH="/sbin:/bin:/usr/sbin:/usr/bin:/usr/local/sbin:/usr/local/bin"
pgm="${0##*/}"				# Program basename
progdir="${0%/*}"			# Program directory
workdir=$( realpath ${progdir} )	# realpath dir
cd ${workdir}

[ -r ${workdir}/nubectl ] && rm -f ${workdir}/nubectl
[ -d ${workdir}/src ] && rm -rf ${workdir}/src

# Check go install
if [ -z "$( which go )" ]; then
	echo "error: Go is not installed. Please install go: pkg install -y lang/go"
	exit 1
fi

# Check go version
GOVERS="$( go version | cut -d " " -f 3 )"
if [ -z "${GOVERS}" ]; then
	echo "unable to determine: go version"
	exit 1
fi

export GOPATH="${workdir}"
#export GOPATH="/tmp/nube"
# go get: cannot install cross-compiled binaries when GOBIN is set
#export GOBIN="/tmp/nube"
export GO111MODULE=off

set -e
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
go get || true
# go get: no install location for directory /usr/home/olevole/nubectl outside GOPATH ??
#GOBIN=/tmp go build -ldflags "${LDFLAGS} -extldflags '-static'" -o "${workdir}/nubectl-windows" --no-clean
go build -ldflags "${LDFLAGS} -extldflags '-static'" -o "${workdir}/nubectl-linux"
