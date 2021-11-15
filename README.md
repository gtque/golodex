# golodex
A simple rolodex app with a React front end with a GO API service talking to a GO a backend in front of CouchDB.

This was developed on windows, and run in Ubuntu in WSL. The kubernetes deployment was developed and tested using a local cluster running in Ubuntu using MicroK8s.

# Requirements
* Must have Go installed with all the required compiler(s)
* Must have npm installed
* Have an accessible CouchDB instance. See golodex_infrastructure for deploying CouchDB to a kubernetes cluster.
  * Set an environment variable called `COUCHDB_ADMIN_PW` with the value for the admin account for the couchdb instance.
  * If the admin user for CouchDB is not `admin` you can set another environment variable `COUCHDB_ADMIN_USER` to the correct user name.

## Overview
golodex_api and golodex_data are written in Go. To build them simply cd to their respective directories and run the `./build.sh` script. This will build the plugins as well. Once built, you can run them directly by running the produced artifacts golodex_api and golodex_data respectively.
golodex_ui is a react app. Simply cd to golodex_ui and run `npm install` and then `npm start`. You only have to run `npm_install` the first time.
The UI is accessible from localhost:3000.

Before building you probably ought to go to google and search for Google Client Id and set up your own client id. Then replace all the references that use the client ID to use your new one.
#### https://console.cloud.google.com/home/dashboard

### To deploy to kubernetes
* install docker (at least docker command line)
* cd to the module you want to build (golodex_api, golodex_data, or golodex_ui)
* edit the `publish.sh` script to either change the image report it pushes to or remove that step.
  * This was originally written to publish to a local image repo running in a local kubernetes cluster managed by MicroK8s.
* run `./publish.sh`
* run `./deploy.sh`
  * This will apply the basic kubernetes objects required for the module. These are defined in the kubernetes folder in each module.
  * Check the configmap and secret files to make sure the settings are correct. Especially the Google Client Id.
  * The current Google Client id is my personal client id for testing this app, and while it will work for localhost, it only has three active google accounts allowed so you will have a bad time if you don't use your own.
  * No, I am not going to provide the steps for setting up the goolge client id. There is ample documentation for that online already.

## Setting up GO
#### This is not a place to to see how to formally get a Go development environment set up. The notes below are what I remember doing while trying to get up and running.
### install gcc
* in linux (ubuntu): sudo apt install build-essential
### install protobuff compiler
* sudo apt install protobuf-compiler
* sudo apt install golang-goprotobuf-dev
### install go
* https://golang.org/doc/install
### install additional compilers
* go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
* go install github.com/micro/micro/v3/cmd/protoc-gen-micro
* https://developers.google.com/protocol-buffers/docs/reference/go-generated
### requires git cmd in path so add it if not already set up.
### run gets to pull libraries:
* go get github.com/asim/go-micro/v3/util/log
* go get google.golang.org/grpc
* go get github.com/asim/go-micro/v3
#### you may also need to do:
* export GOPATH=$(pwd)
  * set the go path to the projects root directory
* protoc --proto_path=./server --go_out=./server --micro_out=./server --go_opt=paths=source_relative pages/about/about.proto
  * configure protoc

## Setting up React
* Install npm
* Install react
  * Install google login: npm install react-google-login
  * npm install --save-dev jest-fetch-mock
  * npm i --save-dev enzyme enzyme-adapter-react-16