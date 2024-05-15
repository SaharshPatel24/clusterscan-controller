/*
Copyright 2024.

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

package controller

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/saharshpatel24/clusterscan-controller/api/v1alpha1"
)

// ClusterScanReconciler reconciles a ClusterScan object
type ClusterScanReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func newInt32(val int32) *int32 {
	return &val
}

//+kubebuilder:rbac:groups=api.my.domain,resources=clusterscans,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.my.domain,resources=clusterscans/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.my.domain,resources=clusterscans/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ClusterScan object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *ClusterScanReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	clusterScan := &apiv1alpha1.ClusterScan{}
	err := r.Get(ctx, req.NamespacedName, clusterScan)
	if err != nil {
		log.Error(err, "unable to fetch ClusterScan")
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	for _, deployment := range clusterScan.Spec.Deployment {
		jobSpec := corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    deployment.JobTemplate.Name,
					Image:   deployment.JobTemplate.Image,
					Command: deployment.JobTemplate.Command,
					Args:    deployment.JobTemplate.Args,
				},
			},
			RestartPolicy: corev1.RestartPolicy(deployment.JobTemplate.RestartPolicy),
		}

		toolStatus := apiv1alpha1.ToolStatus{
			Name: deployment.Name,
		}

		if deployment.Reccuring {
			cronJob := &batchv1.CronJob{
				ObjectMeta: metav1.ObjectMeta{
					Name:      deployment.Name,
					Namespace: deployment.Namespace,
				},
				Spec: batchv1.CronJobSpec{
					Schedule: deployment.Schedule,
					JobTemplate: batchv1.JobTemplateSpec{
						Spec: batchv1.JobSpec{
							Template: corev1.PodTemplateSpec{
								Spec: jobSpec,
							},
						},
					},
					SuccessfulJobsHistoryLimit: newInt32(2),
					FailedJobsHistoryLimit:     newInt32(1),
				},
			}
			if err := r.Create(ctx, cronJob); err != nil {
				log.Error(err, "unable to create CronJob for ClusterScan", "CronJob", cronJob)
				toolStatus.Status = "Failed"
				return ctrl.Result{}, err
			}
		} else {
			job := &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      deployment.Name,
					Namespace: deployment.Namespace,
				},
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: jobSpec,
					},
				},
			}
			if err := r.Create(ctx, job); err != nil {
				log.Error(err, "unable to create Job for ClusterScan", "Job", job)
				toolStatus.Status = "Failed"
				return ctrl.Result{}, err
			}
		}

		toolStatus.Status = "Deployed"

		// Handling Status update
		clusterScan.Status.ToolStatuses = append(clusterScan.Status.ToolStatuses, toolStatus)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterScanReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.ClusterScan{}).
		Complete(r)
}
