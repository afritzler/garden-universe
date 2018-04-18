# garden-universe
![garden universe logo](images/logo.png)
---

# Overview
Garden universe renders a Kubernetes landscape which is setup and managed by the [Gardener Project](https://github.com/gardener/gardener) into a dynamic 3D graph. An example landscape visualization can be found [here](images/universe.png).

# Features
* 3D landscape overview
* Seed <-> shoot relationship 
* Identify shoot clusters with errors
* Show latest error message for a shoot

# Installation and Usage

Get the `garden-universe` binary
```
go get github.com/afritzler/garden-universe
```
Start `garden-universe` locally
```
garden-universe serve --kubeconfig=PATH_TO_MY_GARDENER_CLUSTER_KUBECONFIG
```
The web UI can be accessed via http://localhost:3000 and the rendered graph under http://localhost:3000/graph.
With the `--port` flag you can also specify under which port the garden universe server should be exposed (default is 3000).

# Development

To build and run the garden universe
```
git clone https://github.com/afritzler/garden-universe $GOPATH/src/github.com/afritzler/garden-universe
cd $GOPATH/src/github.com/afritzler/garden-universe
go run *.go serve --kubeconfig=PATH_TO_MY_GARDENER_CLUSTER_KUBECONFIG
```
or to build and run it using the executable
```
make
./garden-universe serve --kubeconfig=PATH_TO_MY_GARDENER_CLUSTER_KUBECONFIG
```
To build the Docker image
```
cd $GOPATH/src/github.com/afritzler/garden-universe
make docker-build
```
# Acknowledgements
Garden universe is using [3d-force-graph](https://github.com/vasturiano/3d-force-graph) for rendering. 
