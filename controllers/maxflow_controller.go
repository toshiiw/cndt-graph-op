/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"math"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	graphv1alpha1 "github.com/toshiiw/cndt-graph-op/api/v1alpha1"
)

// MaxFlowReconciler reconciles a MaxFlow object
type MaxFlowReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type GVertex struct {
	Name  string
	Edges [](*graphv1alpha1.DiEdge)
}

//+kubebuilder:rbac:groups=graph.example.valinux.co.jp,resources=maxflows,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=graph.example.valinux.co.jp,resources=maxflows/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=graph.example.valinux.co.jp,resources=maxflows/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *MaxFlowReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var maxflow graphv1alpha1.MaxFlow
	if err := r.Get(ctx, req.NamespacedName, &maxflow); err != nil {
		log.Error(err, "unable to fetch")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var edges graphv1alpha1.DiEdgeList
	if err := r.List(ctx, &edges, client.InNamespace(req.Namespace)); err != nil {
		log.Error(err, "cannot list edges")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	vertices := make(map[string]*GVertex)
	for i, e := range edges.Items {
		for _, vname := range []string{e.Spec.From, e.Spec.To} {
			v, ok := vertices[vname]
			if !ok {
				v = &GVertex{Name: vname, Edges: make([]*graphv1alpha1.DiEdge, 0, 5)}
				vertices[vname] = v
			}
			(*v).Edges = append((*v).Edges, &edges.Items[i])
		}
	}
	log.Info("loaded", "vertices", len(vertices))
	log.V(2).Info("loaded", "vertices", vertices)

	reset_flows := false
	for vname, v := range vertices {
		// this isn't going to work with multiple maxflow problems
		// in a namespace
		if vname == maxflow.Spec.To {
			continue
		}

		var flow int32
		for _, e := range v.Edges {
			if e.Spec.From == vname {
				flow += e.Spec.Allocated
			} else if e.Spec.To == vname {
				flow -= e.Spec.Allocated
			} else {
				log.Error(nil, "unexpected data inconsistency", "vertex", vname, "edges", v.Edges)
				// XXX
				return ctrl.Result{}, nil
			}
		}
		if vname == maxflow.Spec.From {
			if maxflow.Status.Flow != flow {
				log.Info("unmatched flow detected", "vertex", vname, "flow", flow, "total", maxflow.Status.Flow)
				reset_flows = true
			}
		} else if flow != 0 {
			log.Info("unmatched flow detected", "vertex", vname, "flow", flow)
			reset_flows = true
		}
	}
	if reset_flows {
		log.Info("resetting flows", "maxflow", maxflow)
		maxflow.Status.Stale = true
		maxflow.Status.Flow = 0
		var err error
		if err = r.Status().Update(ctx, &maxflow); err != nil {
			log.Error(err, "cannot set stale")
		}
		for _, e := range edges.Items {
			e.Spec.Allocated = 0
			err := r.Update(ctx, &e)
			if err == nil {
				continue
			}
			if strings.Contains(err.Error(), "the object has been modified") {
				var e2 graphv1alpha1.DiEdge
				if err = r.Get(ctx, client.ObjectKey{Namespace: e.Namespace, Name: e.Name}, &e2); err != nil {
					log.Error(err, "cannot refetch DiEdge for reset")
				} else {
					e2.Spec.Allocated = 0
					err = r.Update(ctx, &e2)
				}
			}
			if err != nil {
				log.Error(err, "cannot reset flow", "diedge", e)
			}
		}
		return ctrl.Result{Requeue: true}, err
	}
	if maxflow.Status.Stale {
		maxflow.Status.Stale = false
		if err := r.Status().Update(ctx, &maxflow); err != nil {
			log.Error(err, "cannot clear stale")
		}
	}
	q := make([]*GVertex, 1)
	qi := 0
	wt := make(map[string]int)
	wp := make(map[string](*graphv1alpha1.DiEdge))
	q[0] = vertices[maxflow.Spec.From]
	wt[maxflow.Spec.From] = 0
	for qi < len(q) {
		for _, e := range q[qi].Edges {
			var other string
			if q[qi].Name == e.Spec.From && e.Spec.Allocated < e.Spec.Capacity {
				other = e.Spec.To
			} else if q[qi].Name == e.Spec.To && e.Spec.Allocated > 0 {
				other = e.Spec.From
			} else {
				continue
			}
			_, ok := wt[other]
			if !ok {
				wt[other] = wt[q[qi].Name] + 1
				wp[other] = e
				q = append(q, vertices[other])
			}
		}
		qi += 1
	}
	log.Info("bfs done", "wt", wt)
	if _, ok := wt[maxflow.Spec.To]; !ok {
		if err := r.Get(ctx, req.NamespacedName, &maxflow); err != nil {
			log.Error(err, "unable to fetch")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		if maxflow.Status.Stale {
			log.Info("got stale. requeing", "maxflow", maxflow)
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, nil
	}
	log.Info("augmenting path found", "length", wt[maxflow.Spec.To])
	type PathComp struct {
		Edge    *graphv1alpha1.DiEdge
		reverse bool
	}
	apath := make([]PathComp, wt[maxflow.Spec.To])
	vname := maxflow.Spec.To
	var acap int32
	acap = math.MaxInt32
	i := 0
	for vname != maxflow.Spec.From {
		var ac1 int32
		apath[i].Edge = wp[vname]
		if vname == wp[vname].Spec.To {
			ac1 = wp[vname].Spec.Capacity - wp[vname].Spec.Allocated
			vname = wp[vname].Spec.From
		} else {
			ac1 = wp[vname].Spec.Allocated
			vname = wp[vname].Spec.To
			apath[i].reverse = true
		}
		if ac1 < acap {
			acap = ac1
		}
		i += 1
	}
	maxflow.Status.Flow += acap
	maxflow.Status.Stale = true
	if err := r.Status().Update(ctx, &maxflow); err != nil {
		log.Error(err, "cannot update maxflow", "maxflow", maxflow)
	}
	for i, _ := range apath {
		if apath[i].reverse {
			apath[i].Edge.Spec.Allocated -= acap
		} else {
			apath[i].Edge.Spec.Allocated += acap
		}
		if apath[i].Edge.Spec.Allocated < 0 || apath[i].Edge.Spec.Allocated > apath[i].Edge.Spec.Capacity {
			err := fmt.Errorf("Unexpected inconsistency")
			log.Error(err, "failed to update flow")
			log.Info("error flow", "edge", apath[i].Edge)
			return ctrl.Result{}, err
		}
		if err := r.Update(ctx, apath[i].Edge); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{Requeue: true}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MaxFlowReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&graphv1alpha1.MaxFlow{}).
		Complete(r)
}
