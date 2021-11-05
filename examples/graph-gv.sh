#!/bin/sh

KIND_HOST=${KIND_HOST:-127.0.0.1}
KUBECTL=${KUBECTL:-kubectl}

T=`mktemp -d`
cd $T
Q=

while true; do
	ssh $KIND_HOST $KUBECTL get edge -o go-template=\''{{range .items}}{{index .spec.vertex 0}}{{" -- "}}{{index .spec.vertex 1}}{{"[label=\""}}{{.spec.weight}}{{"\""}}{{if .status.inmaxcut}}{{" color=brown"}}{{end}}{{"]\n"}}{{end}}'\' > status
	if diff -q status status.1 > /dev/null; then
		:
	else
		echo "graph G {
		{ node [style=filled]" > graph.gv
		ssh $KIND_HOST $KUBECTL get maxcut problem -o go-template=\''{{range $i, $e := .status.vertexset}}{{$e}}{{" [fillcolor=yellow]\n"}}{{end}}'\' >> graph.gv
		echo "}" >> graph.gv
		cat status >> graph.gv
		echo "}" >> graph.gv
		dot -Tpng graph.gv > graph.png
		mv status status.1
		if [ -z "$Q" ]; then
			qiv -eT graph.png&
			Q=1
		fi
	fi
	sleep 2
done
