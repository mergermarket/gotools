#!/bin/bash -e

function lintCode() {
	golangci-lint run --timeout=10m
}

function fmtCode() {
	go fmt $(go list ./... | grep -v /vendor/)
}

function runTests() {
    # We need to do a bit of fiddling to generate coverage data from multiple packages and merge them.
    (
        if [ ! -d coverage ]; then
            mkdir coverage
        fi
        echo "mode: count" > coverage/coverage-all.out
        for pkg in $(go list ./... | grep -v /vendor/)
        do
          echo "pkg=$pkg"
          go test -v -coverprofile=coverage/coverage.out -covermode=count $pkg
          if [ -f coverage/coverage.out ]; then
            tail -n +2 coverage/coverage.out >> coverage/coverage-all.out
          fi
        done
        go test  -v -tags=integration -coverprofile=coverage/coverage.out -covermode=count
        tail -n +2 coverage/coverage.out >> coverage/coverage-all.out
        go tool cover -html=coverage/coverage-all.out -o coverage/coverage.html
    )
}

function buildCode() {
    go build -v
}

