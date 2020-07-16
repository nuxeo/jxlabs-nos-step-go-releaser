#!/bin/bash

echo "promoting the new version ${VERSION} to downstream repositories"

jx step create pr go --name github.com/nuxeo/jxlabs-nos-step-go-releaser --version ${VERSION} --build "make build" --repo https://github.com/nuxeo/jxlabs-nos-jxl.git
