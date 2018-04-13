# garden-universe
![garden universe logo](images/logo.png)
---

# Overview
Garden universe renders a Kubernetes landscape which is setup and managed by the [Gardener Project](https://github.com/gardener/gardener) into a dynamic 3D graph. An example landscape redering can be found [here](images/universe.png).

# Development

To locally run the garden universe
```
git clone https://github.com/afritzler/garden-universe $GOPATH/src/github.com/afritzler/garden-universe
cd $GOPATH/src/github.com/afritzler/garden-universe
go run *.go serve --kubeconfig=PATH_TO_MY_GARDEN_CLUSTER_KUBECONFIG
```
The web UI can be accessed via http://localhost:3000

To build the executable run
```
cd $GOPATH/src/github.com/afritzler/garden-universe
make
```

To build the Docker image
```
cd $GOPATH/src/github.com/afritzler/garden-universe
make docker-build
```

# Acknowledgements
Gardener universe is using the [3d-force-graph](https://github.com/vasturiano/3d-force-graph) for rendering. 
