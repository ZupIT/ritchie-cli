#!/bin/sh

cd /home/application

./wait-for-it.sh "stubby4j:8882" && echo "stubby4j is up"

export REPO_URL=http://stubby4j:8882/formulas

PACKAGE_NAMES=$(go list ./pkg/... | grep -v vendor/)

gotestsum --format=short-verbose --junitfile $TEST_RESULTS_DIR/gotestsum-report.xml -- -p 2 -cover -coverprofile=coverage.txt $PACKAGE_NAMES

testStatus=$?
if [ $testStatus -ne 0 ]; then
    echo "Tests failed"
    exit 1
fi
