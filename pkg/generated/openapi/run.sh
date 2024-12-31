#!/bin/bash

cd ../../..

./hack/update-codegen.sh openapi

cd -

go test -v -test.run TestWriteSchema
