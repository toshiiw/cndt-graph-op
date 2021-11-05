# A graph optimization solver operator

This is a toy example of the k8s operator framework.
It can solve max flow and max cut problems defined as k8s CRDs.

## Running

Refer to the operator SDK documentation for details, but can be as simple as:

```
$ make
$ make manifests
$ make install
$ make run
```

## Examples

These are located under the [examples](./examples/) directory.

Graph definitions and problems:

- examples/digraph.yaml
- examples/graph.yaml

Random graph yaml generator:

- examples/graph-gen.py

Graph visualizing scripts:

- examples/graph-gv.sh
- examples/digraph-gv.sh
