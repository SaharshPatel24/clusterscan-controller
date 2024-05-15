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
	"k8s.io/apimachinery/pkg/types"
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

type JobHandleResult struct {
	Result ctrl.Result
	Error  error
	Status string
}

func newInt32(val int32) *int32 {
	return &val
}

func (r *ClusterScanReconciler) handleJob(ctx context.Context, deployment apiv1alpha1.DeploymentSpec) JobHandleResult {
	// Define the PodSpec
	log := log.FromContext(ctx)
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

	var job batchv1.Job
	jobName := types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}
	jobExists := true
	if err := r.Get(ctx, jobName, &job); err != nil {
		if errors.IsNotFound(err) {
			jobExists = false
		} else {
			log.Error(err, "unable to fetch CronJob")
			return JobHandleResult{
				Result: ctrl.Result{},
				Error:  err,
				Status: "Unknow Error",
			}
		}
	}

	if !jobExists {
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
			return JobHandleResult{
				Result: ctrl.Result{},
				Error:  err,
				Status: "Failed creating job",
			}
		}

		return JobHandleResult{
			Result: ctrl.Result{},
			Error:  nil,
			Status: "Successfully Deployed",
		}
	}

	return JobHandleResult{
		Result: ctrl.Result{},
		Error:  nil,
		Status: "Currently Job is running",
	}
}

func (r *ClusterScanReconciler) handleCronJob(ctx context.Context, deployment apiv1alpha1.DeploymentSpec) JobHandleResult {
	// Fetch the CronJob
	log := log.FromContext(ctx)
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

	var cronJob batchv1.CronJob
	cronJobName := types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}
	cronJobExists := true
	if err := r.Get(ctx, cronJobName, &cronJob); err != nil {
		if errors.IsNotFound(err) {
			cronJobExists = false
		} else {
			log.Error(err, "unable to fetch CronJob")
			return JobHandleResult{
				Result: ctrl.Result{},
				Error:  err,
				Status: "Unknow Error",
			}
		}
	}

	// If the CronJob doesn't exist, create it
	if !cronJobExists {
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
			return JobHandleResult{
				Result: ctrl.Result{},
				Error:  err,
				Status: "Failed creating job",
			}
		}

		return JobHandleResult{
			Result: ctrl.Result{},
			Error:  nil,
			Status: "Successfully Deployed",
		}
	}

	return JobHandleResult{
		Result: ctrl.Result{},
		Error:  nil,
		Status: "Currently CronJob Successfully Running",
	}
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
	clusterScan.Status.ToolStatuses = nil
	for _, deployment := range clusterScan.Spec.Deployment {
		var result JobHandleResult

		if deployment.Reccuring {
			result = r.handleCronJob(ctx, deployment)
		} else {
			result = r.handleJob(ctx, deployment)
		}

		if result.Error != nil {
			log.Error(result.Error, "error handling job or cronjob")
			return result.Result, result.Error
		}

		toolstatus := apiv1alpha1.ToolStatus{
			Name:   deployment.JobTemplate.Name,
			Status: result.Status,
		}

		// Update the status of the custom resource
		clusterScan.Status.ToolStatuses = append(clusterScan.Status.ToolStatuses, toolstatus)
	}

	if err := r.Status().Update(ctx, clusterScan); err != nil {
		log.Error(err, "unable to update ClusterScan status")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterScanReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.ClusterScan{}).
		Complete(r)
}
