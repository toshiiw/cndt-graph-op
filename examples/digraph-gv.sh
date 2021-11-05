#!/bin/sh

KIND_HOST=${KIND_HOST:-127.0.0.1}
KUBECTL=${KUBECTL:-kubectl}

T=`mktemp -d`
cd $T
Q=

while true; do
	ssh $KIND_HOST $KUBECTL get diedge -o go-template=\''{{range .items}}{{.spec.from}}{{" -> "}}{{.spec.to}}{{"\t[label=\""}}{{.spec.allocated}}{{"/"}}{{.spec.capacity}}{{"\""}}{{if eq .spec.capacity .spec.allocated}}{{" color=red"}}{{end}}{{"]\n"}}{{end}}'\' > status

	if diff -q status status.1 > /dev/null; then
		:
	else
		echo "digraph G {" > graph.gv
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
