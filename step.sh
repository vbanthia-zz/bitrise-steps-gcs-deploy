#!/bin/bash
set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPT_DIR}"

go get -u cloud.google.com/go/storage

go get -u -d github.com/bitrise-io/go-utils/log
go get -u -d github.com/bitrise-io/go-utils/pathutil

go get -u -d google.golang.org/api/iterator
go get -u -d google.golang.org/api/option

go run ./gcs_upload.go ${project_name}
