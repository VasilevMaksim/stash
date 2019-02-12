package controller

import (
	"fmt"

	"github.com/appscode/go/log"
	workload_api "github.com/appscode/kubernetes-webhook-util/apis/workload/v1"
	batch_util "github.com/appscode/kutil/batch/v1beta1"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/queue"
	api_v1beta1 "github.com/appscode/stash/apis/stash/v1beta1"
	"github.com/appscode/stash/client/clientset/versioned/scheme"
	"github.com/appscode/stash/pkg/docker"
	"github.com/appscode/stash/pkg/eventer"
	"github.com/appscode/stash/pkg/util"
	"github.com/golang/glog"
	"k8s.io/api/batch/v1beta1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/reference"
)

func (c *StashController) initBackupConfigurationWatcher() {
	c.bupcInformer = c.stashInformerFactory.Stash().V1beta1().BackupConfigurations().Informer()
	c.bupcQueue = queue.New("BackupConfiguration", c.MaxNumRequeues, c.NumThreads, c.runBackupConfigurationInjector)
	c.bupcInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if b, ok := obj.(*api_v1beta1.BackupConfiguration); ok {
				ref, rerr := reference.GetReference(scheme.Scheme, b)
				if rerr == nil {
					c.recorder.Eventf(
						ref,
						core.EventTypeWarning,
						eventer.EventReasonInvalidBackupConfiguration,
						"Reason: %v",
						//err,
					)
				}
				//return
			}

			queue.Enqueue(c.bupcQueue.GetQueue(), obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldBupc, ok := oldObj.(*api_v1beta1.BackupConfiguration)
			if !ok {
				log.Errorln("Invalid BackupConfiguration object")
				return
			}
			newbupc, ok := newObj.(*api_v1beta1.BackupConfiguration)
			if !ok {
				log.Errorln("Invalid BackupConfiguration object")
				return
			}
			ref, rerr := reference.GetReference(scheme.Scheme, newbupc)
			if rerr == nil {
				c.recorder.Eventf(
					ref,
					core.EventTypeWarning,
					eventer.EventReasonInvalidBackupConfiguration,
					"Reason: %v",
					//err,
				)
			}
			if !util.BackupConfigurationEqual(oldBupc, newbupc) {
				queue.Enqueue(c.bupcQueue.GetQueue(), newObj)
			}

		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(c.bupcQueue.GetQueue(), obj)

		},
	})
	c.bupcLister = c.stashInformerFactory.Stash().V1beta1().BackupConfigurations().Lister()
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the deployment to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (c *StashController) runBackupConfigurationInjector(key string) error {
	obj, exists, err := c.bupcInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}
	if !exists {
		glog.Errorf("BackupConfiguration %s does not exit anymore\n", key)
		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}

		c.EnsureSidecarDeleted2(namespace, name)
		err = c.EnsureCronJobDeleted(namespace, name)
		fmt.Println(err)
		if err != nil {
			return err
		}
	} else {
		backupconfiguration := obj.(*api_v1beta1.BackupConfiguration)
		glog.Info("Sync/Add/Update for BackupConfiguration %s", backupconfiguration.GetName())
		if backupconfiguration.Spec.Target != nil {
			if backupconfiguration.Spec.Target.Ref != nil {
				kind := backupconfiguration.Spec.Target.Ref.Kind
				if kind == workload_api.KindDeployment || kind == workload_api.KindDaemonSet || kind == workload_api.KindReplicationController || kind == workload_api.KindReplicaSet || kind == workload_api.KindStatefulSet {
					c.EnsureSidecar2(backupconfiguration)
				}

			}
			err := c.EnsureCronJob(backupconfiguration)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (c *StashController) EnsureSidecarDeleted2(namespace, name string) {
	if resources, err := c.dpLister.Deployments(namespace).List(labels.Everything()); err == nil {
		for _, resource := range resources {
			backupconfiguration, err := util.GetAppliedBackupConfiguration(resource.Annotations)
			if err != nil {
				if ref, e2 := reference.GetReference(scheme.Scheme, resource); e2 == nil {
					c.recorder.Eventf(
						ref,
						core.EventTypeWarning,
						eventer.EventReasonInvalidBackupConfiguration,
						"Reason: %s",
						err.Error(),
					)
				}
			} else if backupconfiguration != nil && backupconfiguration.Namespace == namespace && backupconfiguration.Name == name {
				key, err := cache.MetaNamespaceKeyFunc(resource)
				if err == nil {
					c.dpQueue.GetQueue().Add(key)
				}

			}

		}
	}
	if resources, err := c.dsLister.DaemonSets(namespace).List(labels.Everything()); err == nil {
		for _, resource := range resources {
			backupconfiguration, err := util.GetAppliedBackupConfiguration(resource.Annotations)
			if err != nil {
				if ref, e2 := reference.GetReference(scheme.Scheme, resource); e2 == nil {
					c.recorder.Eventf(
						ref,
						core.EventTypeWarning,
						eventer.EventReasonInvalidBackupConfiguration,
						"Reason: %s",
						err.Error(),
					)
				}
			} else if backupconfiguration != nil && backupconfiguration.Namespace == namespace && backupconfiguration.Name == name {
				key, err := cache.MetaNamespaceKeyFunc(resource)
				if err != nil {
					c.dsQueue.GetQueue().Add(key)
				}
			}
		}
	}
	if resources, err := c.ssLister.StatefulSets(namespace).List(labels.Everything()); err == nil {
		for _, resource := range resources {
			backupconfiguration, err := util.GetAppliedBackupConfiguration(resource.Annotations)
			if err != nil {
				if ref, e2 := reference.GetReference(scheme.Scheme, resource); e2 == nil {
					c.recorder.Eventf(
						ref,
						core.EventTypeWarning,
						eventer.EventReasonInvalidBackupConfiguration,
						"Reason: %s",
						err.Error(),
					)
				}
			} else if backupconfiguration != nil && backupconfiguration.Namespace == namespace && backupconfiguration.Name == name {
				key, err := cache.MetaNamespaceKeyFunc(resource)
				if err != nil {
					c.ssQueue.GetQueue().Add(key)
				}
			}
		}
	}
	if resources, err := c.rcLister.ReplicationControllers(namespace).List(labels.Everything()); err == nil {
		for _, resource := range resources {
			backupconfiguration, err := util.GetAppliedBackupConfiguration(resource.Annotations)
			if err != nil {
				if ref, e2 := reference.GetReference(scheme.Scheme, resource); e2 == nil {
					c.recorder.Eventf(
						ref,
						core.EventTypeWarning,
						eventer.EventReasonInvalidBackupConfiguration,
						"Reason: %s",
						err.Error(),
					)
				}
			} else if backupconfiguration != nil && backupconfiguration.Namespace == namespace && backupconfiguration.Name == name {
				key, err := cache.MetaNamespaceKeyFunc(resource)
				if err != nil {
					c.rcQueue.GetQueue().Add(key)
				}
			}
		}
	}
	if resources, err := c.rsLister.ReplicaSets(namespace).List(labels.Everything()); err == nil {
		for _, resource := range resources {
			backupconfiguration, err := util.GetAppliedBackupConfiguration(resource.Annotations)
			if err != nil {
				if ref, e2 := reference.GetReference(scheme.Scheme, resource); e2 == nil {
					c.recorder.Eventf(
						ref,
						core.EventTypeWarning,
						eventer.EventReasonInvalidBackupConfiguration,
						"Reason: %s",
						err.Error(),
					)
				}
			} else if backupconfiguration != nil && backupconfiguration.Namespace == namespace && backupconfiguration.Name == name {
				key, err := cache.MetaNamespaceKeyFunc(resource)
				if err != nil {
					c.rsQueue.GetQueue().Add(key)
				}
			}
		}
	}
}

func (c *StashController) EnsureSidecar2(backupconfiguration *api_v1beta1.BackupConfiguration) {
	resource_name := backupconfiguration.Spec.Target.Ref.Name
	switch backupconfiguration.Spec.Target.Ref.Kind {
	case workload_api.KindDeployment:
		if resource, err := c.dpLister.Deployments(backupconfiguration.Namespace).Get(resource_name); err == nil {
			key, err := cache.MetaNamespaceKeyFunc(resource)
			if err == nil {
				c.dpQueue.GetQueue().Add(key)
			}
		}
	case workload_api.KindDaemonSet:
		if resource, err := c.dsLister.DaemonSets(backupconfiguration.Namespace).Get(resource_name); err == nil {
			key, err := cache.MetaNamespaceKeyFunc(resource)
			if err == nil {
				c.dsQueue.GetQueue().Add(key)
			}
		}
	case workload_api.KindStatefulSet:
		if resource, err := c.ssLister.StatefulSets(backupconfiguration.Namespace).Get(resource_name); err == nil {
			key, err := cache.MetaNamespaceKeyFunc(resource)
			if err == nil {
				c.ssQueue.GetQueue().Add(key)
			}
		}
	case workload_api.KindReplicationController:
		if resource, err := c.rcLister.ReplicationControllers(backupconfiguration.Namespace).Get(resource_name); err == nil {
			key, err := cache.MetaNamespaceKeyFunc(resource)
			if err == nil {
				c.rcQueue.GetQueue().Add(key)
			}
		}
	case workload_api.KindReplicaSet:
		if resource, err := c.rsLister.ReplicaSets(backupconfiguration.Namespace).Get(resource_name); err == nil {
			key, err := cache.MetaNamespaceKeyFunc(resource)
			if err == nil {
				c.rsQueue.GetQueue().Add(key)
			}
		}

	}

}

func (c *StashController) EnsureCronJob(backupconfiguration *api_v1beta1.BackupConfiguration) error {
	image := docker.Docker{
		Registry: c.DockerRegistry,
		Image:    docker.ImageStash,
		Tag:      c.StashImageTag,
	}

	meta := metav1.ObjectMeta{
		Name:      backupconfiguration.Name,
		Namespace: backupconfiguration.Namespace,
	}

	if c.EnableRBAC {
		ref, err := reference.GetReference(scheme.Scheme, backupconfiguration)
		if err != nil {
			return err
		}
		if err := c.ensureCronJobRBAC(ref); err != nil {
			return fmt.Errorf("error ensuring rbac for kubectl cron job %s, reason: %s", meta.Name, err)
		}
	}

	_, _, err := batch_util.CreateOrPatchCronJob(c.kubeClient, meta, func(in *v1beta1.CronJob) *v1beta1.CronJob {
		// set backup-configuration as cron-job owner
		in.OwnerReferences = []metav1.OwnerReference{
			{
				APIVersion: api_v1beta1.SchemeGroupVersion.String(),
				Kind:       api_v1beta1.ResourceKindBackupConfiguration,
				Name:       backupconfiguration.Name,
				UID:        backupconfiguration.UID,
			},
		}
		in.Spec.Schedule = backupconfiguration.Spec.Schedule
		if in.Spec.JobTemplate.Labels == nil {
			in.Spec.JobTemplate.Labels = map[string]string{}
		}
		in.Spec.JobTemplate.Labels["app"] = util.AppLabelStash
		in.Spec.JobTemplate.Spec.Template.Spec.Containers = core_util.UpsertContainer(
			in.Spec.JobTemplate.Spec.Template.Spec.Containers,
			core.Container{
				Name:            backupconfiguration.Name,
				ImagePullPolicy: "Always",
				Image:           image.ToContainerImage(),
				Args: []string{
					"backup_session",
					fmt.Sprintf("--backupSessionName=%s", backupconfiguration.Name),
					fmt.Sprintf("--backupSessionNamespace=%s", backupconfiguration.Namespace),
				},
			})
		in.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy = core.RestartPolicyNever
		in.Spec.JobTemplate.Spec.Template.Spec.ServiceAccountName = in.Name
		return in

	})
	if err != nil {
		return err
	}

	return nil
}

func (c *StashController) EnsureCronJobDeleted(namespace, name string) error {
	//backupconfiguration, err := c.bupcLister.BackupConfigurations(namespace).Get(name)
	//ref, err := reference.GetReference(scheme.Scheme, backupconfiguration)
	//if err != nil {
	//	return err
	//}
	//meta := metav1.ObjectMeta{
	//	Name:      ref.Name,
	//	Namespace: ref.Namespace,
	//}
	//core_util.EnsureOwnerReference(&meta, ref)
	deletePolicy := metav1.DeletePropagationBackground
	if err := c.kubeClient.BatchV1beta1().CronJobs(namespace).Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	return nil
}
