#!/bin/sh

PACKAGE_NAMES=$(go list ./pkg/... | grep -v vendor/)

go mod download
gotestsum --format=short-verbose --junitfile $TEST_RESULTS_DIR/gotestsum-report.xml -- -p 2 -cover -coverprofile=coverage.txt $PACKAGE_NAMES

testStatus=$?
if [ $testStatus -ne 0 ]; then
    echo "Tests failed"
    exit 1
fi
