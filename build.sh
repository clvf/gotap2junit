#!/bin/env bash

if ! command -v go > /dev/null; then
    echo "go command cannot be found!" >&2
    exit 1
fi

set -e

scriptdir="$(cd "$(dirname $0)" && pwd)"

mkdir -p "${scriptdir}"/bin && go build -o "${scriptdir}"/bin ./...

echo "${scriptdir}"/bin/*
