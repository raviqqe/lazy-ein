#!/bin/sh

set -e

export CGO_CPPFLAGS="$(llvm-config --cppflags)"
export CGO_CXXFLAGS=-std=c++11
export CGO_LDFLAGS="$(llvm-config --ldflags --libs --system-libs all)"

go build -tags byollvm
