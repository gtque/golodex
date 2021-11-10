#!/bin/bash

echo "building app image"
docker build --tag localhost:32000/golodex-ui:dev .
echo "build finished"
echo "pushing to local image repo"
docker push localhost:32000/golodex-ui:dev
echo "pushed"