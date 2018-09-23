package hellostateful

import (
	"fmt"

	"path/filepath"

	"github.com/joatmon08/hello-stateful-operator/pkg/apis/hello-stateful/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	RestoreName = "restore-%s"
)

// Restore generates a CronJob that backs up
// the backend of the StatefulSet.
func Restore(hs *v1alpha1.HelloStateful) error {
	if !hs.Spec.RestoreFromExisting || hs.Status.IsRestored {
		return nil
	}
	restoreJob, err := newJob(hs)
	if err != nil {
		return err
	}
	err = sdk.Create(restoreJob)
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create restore job: %v", err)
		return err
	}
	err = sdk.Get(restoreJob)
	if err != nil {
		logrus.Errorf("Failed to get restore job: %v", err)
		return err
	}
	hs.Status.IsRestored = true
	err = sdk.Update(hs)
	if err != nil {
		return fmt.Errorf("failed to update hellostateful status: %v", err)
	}
	return nil
}

func newJob(cr *v1alpha1.HelloStateful) (*batchv1.Job, error) {
	labels := labelsForHelloStatefulRestore(cr.ObjectMeta.Name)
	backupHostPathType := corev1.HostPathDirectory
	appHostPathType := corev1.HostPathDirectory
	hostpathDirectory := filepath.Dir(cr.Status.BackendVolumes[0])
	volumeDirectory := filepath.Base(cr.Status.BackendVolumes[0])
	logrus.Infof("Restoring to: %s", hostpathDirectory)
	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(RestoreName, cr.ObjectMeta.Name),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						corev1.Container{
							Name:            BACKUPCONTAINERNAME,
							Image:           BACKUPIMAGE,
							ImagePullPolicy: IMAGEPULLPOLICY,
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      BACKUPVOLUME,
									MountPath: BACKUPFOLDER,
								},
								corev1.VolumeMount{
									Name:      APPVOLUME,
									MountPath: hostpathDirectory,
								},
							},
							Args: []string{
								fmt.Sprintf("%s/%s", BACKUPFOLDER, volumeDirectory),
								hostpathDirectory,
							},
						},
					},
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: BACKUPVOLUME,
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: HOSTBACKUPFOLDER,
									Type: &backupHostPathType,
								},
							},
						},
						corev1.Volume{
							Name: APPVOLUME,
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: hostpathDirectory,
									Type: &appHostPathType,
								},
							},
						},
					},
				},
			},
		},
	}
	addOwnerRefToObject(job, asOwner(cr))
	return job, nil
}
