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
	"math/rand"
	"sort"

	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	graphv1alpha1 "github.com/toshiiw/cndt-graph-op/api/v1alpha1"
)

// MaxCutReconciler reconciles a MaxCut object
type MaxCutReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=graph.example.valinux.co.jp,resources=maxcuts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=graph.example.valinux.co.jp,resources=maxcuts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=graph.example.valinux.co.jp,resources=maxcuts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MaxCut object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *MaxCutReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var maxcut graphv1alpha1.MaxCut
	if err := r.Get(ctx, req.NamespacedName, &maxcut); err != nil {
		log.Error(err, "unable to fetch maxcut")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var edges graphv1alpha1.EdgeList
	if err := r.List(ctx, &edges, client.InNamespace(req.Namespace)); err != nil {
		log.Error(err, "unable to list edges")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	vertices := make(map[string]int)
	for _, e := range edges.Items {
		for _, v := range e.Spec.Vertex {
			vertices[v] = 0
		}
	}

	var i int32
	for i = 0; i < maxcut.Spec.Iteration; i++ {
		for k, _ := range vertices {
			vertices[k] = rand.Intn(2)
		}
		total := 0
		for _, e := range edges.Items {
			if vertices[e.Spec.Vertex[0]] != vertices[e.Spec.Vertex[1]] {
				total += int(e.Spec.Weight)
			}
		}
		if total > int(maxcut.Status.Total) {
			m := [][]string{make([]string, 0), make([]string, 0)}
			for v, vi := range vertices {
				m[vi] = append(m[vi], v)
			}
			for j := 0; j < 2; j++ {
				sort.Strings(m[j])
			}
			if m[0][0] < m[1][0] {
				maxcut.Status.VertexSet = m[0]
				maxcut.Status.ComplementVertexSet = m[1]
			} else {
				maxcut.Status.VertexSet = m[1]
				maxcut.Status.ComplementVertexSet = m[0]
			}
			maxcut.Status.Total = int32(total)
			if err := r.Status().Update(ctx, &maxcut); err != nil {
				log.Error(err, "failed to update status")
				return ctrl.Result{}, err
			}

			for _, e := range edges.Items {
				newstatus := vertices[e.Spec.Vertex[0]] != vertices[e.Spec.Vertex[1]]
				if newstatus != e.Status.InMaxCut {
					e.Status.InMaxCut = newstatus
					if err := r.Status().Update(ctx, &e); err != nil {
						log.Error(err, "failed to update status", "edge", e)
					}
				}
			}
			return ctrl.Result{Requeue: true}, nil
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MaxCutReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&graphv1alpha1.MaxCut{}).
		Complete(r)
}
