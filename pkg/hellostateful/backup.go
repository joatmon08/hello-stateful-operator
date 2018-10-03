package hellostateful

import (
	"fmt"

	"github.com/joatmon08/hello-stateful-operator/pkg/apis/hello-stateful/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Constants for the backup job
const (
	BackupVolumeName    = "backup"
	BackupFolder        = "/tmp/backup"
	HostBackupFolder    = "/tmp/backup"
	BackupImage         = "joatmon08/hello-stateful-backup:latest"
	BackupContainerName = "hello-stateful-backup"
)

// Backup generates a CronJob that backs up
// the backend of the StatefulSet.
func Backup(hs *v1alpha1.HelloStateful) error {
	if len(hs.Status.BackendVolumes) < 1 {
		return fmt.Errorf("volumes not created yet, did not create backup job")
	}
	backupJob, err := newCronJob(hs)
	if err != nil {
		return err
	}
	err = sdk.Create(backupJob)
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create backup job: %v", err)
		return err
	}
	err = sdk.Get(backupJob)
	if err != nil {
		logrus.Errorf("Failed to get backup job: %v", err)
		return err
	}
	return nil
}

func newCronJob(cr *v1alpha1.HelloStateful) (*batchv1beta1.CronJob, error) {
	labels := labelsForHelloStatefulBackup(cr.ObjectMeta.Name)
	backupHostPathType := corev1.HostPathDirectoryOrCreate
	appHostPathType := corev1.HostPathDirectory
	job := &batchv1beta1.CronJob{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CronJob",
			APIVersion: "batch/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("backup-%s", cr.ObjectMeta.Name),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: batchv1beta1.CronJobSpec{
			Schedule: cr.Spec.BackupSchedule,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							RestartPolicy: corev1.RestartPolicyOnFailure,
							Containers: []corev1.Container{
								corev1.Container{
									Name:            BackupContainerName,
									Image:           BackupImage,
									ImagePullPolicy: ImagePullPolicy,
									VolumeMounts: []corev1.VolumeMount{
										corev1.VolumeMount{
											Name:      BackupVolumeName,
											MountPath: BackupFolder,
										},
										corev1.VolumeMount{
											Name:      AppVolumeName,
											MountPath: cr.Status.BackendVolumes[0],
										},
									},
									Args: []string{
										cr.Status.BackendVolumes[0],
										BackupFolder,
									},
								},
							},
							Volumes: []corev1.Volume{
								corev1.Volume{
									Name: BackupVolumeName,
									VolumeSource: corev1.VolumeSource{
										HostPath: &corev1.HostPathVolumeSource{
											Path: HostBackupFolder,
											Type: &backupHostPathType,
										},
									},
								},
								corev1.Volume{
									Name: AppVolumeName,
									VolumeSource: corev1.VolumeSource{
										HostPath: &corev1.HostPathVolumeSource{
											Path: cr.Status.BackendVolumes[0],
											Type: &appHostPathType,
										},
									},
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
