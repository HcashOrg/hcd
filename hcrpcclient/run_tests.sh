#!/usr/bin/env bash
set -ex

# The script does automatic checking on a Go package and its sub-packages,
# including:
# 1. gofmt         (http://golang.org/cmd/gofmt/)
# 2. golint        (https://github.com/golang/lint)
# 3. go vet        (http://golang.org/cmd/vet)
# 4. goimports     (https://github.com/bradfitz/goimports)
# 5. unconvert     (https://github.com/mdempsky/unconvert)

# gometalinter (github.com/alecthomas/gometalinter) is used to run each each
# static checker.

# To run on docker on windows, symlink /mnt/c to /c and then execute the script
# from the repo path under /c.  See:
# https://github.com/Microsoft/BashOnWindows/issues/1854
# for more details.

#Default GOVERSION
GOVERSION=${1:-1.8}
REPO=hcrpcclient

TESTCMD="test -z \"\$(gometalinter --disable-all \
  --enable=gofmt \
  --enable=golint \
  --enable=vet \
  --enable=goimports \
  --enable=unconvert \
  --vendor \
  --deadline=20s ./... | grep -v 'ALL_CAPS\|OP_' 2>&1 | tee /dev/stderr)\""

if [ $GOVERSION == "local" ]; then
    eval $TESTCMD
    exit
fi

DOCKER_IMAGE_TAG=decred-golang-builder-$GOVERSION

docker pull decred/$DOCKER_IMAGE_TAG

docker run --rm -it -v $(pwd):/src decred/$DOCKER_IMAGE_TAG /bin/bash -c "\
  rsync -ra --filter=':- .gitignore'  \
  /src/ /go/src/github.com/hcd/$REPO/ && \
  cd github.com/hcd/$REPO/ && \
  glide install && \
  go install \$(glide novendor) && \
  $TESTCMD
"

echo "------------------------------------------"
echo "Tests complete."
