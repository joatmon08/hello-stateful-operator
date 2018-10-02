package e2eutil

import (
	"testing"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

const BACKUPJOBTHRESHOLDSECONDS = 60

func WaitForJob(t *testing.T, kubeclient kubernetes.Interface, namespace, name string, retryInterval, timeout time.Duration) error {
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		t.Logf("Checking for job %s in namespace %s\n", name, namespace)
		job, err := kubeclient.BatchV1().Jobs(namespace).Get(name, metav1.GetOptions{IncludeUninitialized: true})
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s job\n", name)
				return false, nil
			}
			return false, err
		}
		if job.Status.Succeeded > 0 {
			return true, nil
		}
		t.Logf("Waiting for job %s (%d)\n", name, job.Status.Succeeded)
		return false, nil
	})
	if err != nil {
		return err
	}
	t.Log("Job completed\n")
	return nil
}

func WaitForCronJob(t *testing.T, kubeclient kubernetes.Interface, namespace, name string, retryInterval, timeout time.Duration) error {
	var secondsSinceLastRun int
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		t.Logf("Checking for cronjob %s in namespace %s\n", name, namespace)
		job, err := kubeclient.BatchV1beta1().CronJobs(namespace).Get(name, metav1.GetOptions{IncludeUninitialized: true})
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s job\n", name)
				return false, nil
			}
			return false, err
		}
		if job.Status.LastScheduleTime != nil {
			secondsSinceLastRun = int(time.Now().Sub(job.Status.LastScheduleTime.Time).Seconds())
			if secondsSinceLastRun < BACKUPJOBTHRESHOLDSECONDS {
				return true, nil
			}
		}
		t.Logf("Waiting for cron job %s (%d)\n", name, len(job.Status.Active))
		return false, nil
	})
	if err != nil {
		return err
	}
	t.Logf("CronJob available (%d seconds /%d seconds)\n", secondsSinceLastRun, BACKUPJOBTHRESHOLDSECONDS)
	return nil
}

func WaitForStatefulSet(t *testing.T, kubeclient kubernetes.Interface, namespace, name string, replicas int, retryInterval, timeout time.Duration) error {
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		statefulset, err := kubeclient.AppsV1().StatefulSets(namespace).Get(name, metav1.GetOptions{IncludeUninitialized: true})
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s statefulset\n", name)
				return false, nil
			}
			return false, err
		}

		if int(statefulset.Status.ReadyReplicas) == replicas {
			return true, nil
		}
		t.Logf("Waiting for full availability of %s statefulset (%d/%d)\n", name, statefulset.Status.ReadyReplicas, replicas)
		return false, nil
	})
	if err != nil {
		return err
	}
	t.Logf("StatefulSet available (%d/%d)\n", replicas, replicas)
	return nil
}
