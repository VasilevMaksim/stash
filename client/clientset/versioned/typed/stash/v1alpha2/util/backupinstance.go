package util

import (
	"fmt"

	"github.com/appscode/kutil"
	restic_util "github.com/appscode/kutil/tools/restic"
	"github.com/appscode/stash/apis"
	api "github.com/appscode/stash/apis/stash/v1alpha2"
	cs "github.com/appscode/stash/client/clientset/versioned/typed/stash/v1alpha2"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

func CreateOrPatchBackupInstance(c cs.StashV1alpha2Interface, meta metav1.ObjectMeta, transform func(alert *api.BackupInstance) *api.BackupInstance) (*api.BackupInstance, kutil.VerbType, error) {
	cur, err := c.BackupInstances(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating BackupInstance %s/%s.", meta.Namespace, meta.Name)
		out, err := c.BackupInstances(meta.Namespace).Create(transform(&api.BackupInstance{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BackupInstance",
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchBackupInstance(c, cur, transform)
}

func PatchBackupInstance(c cs.StashV1alpha2Interface, cur *api.BackupInstance, transform func(*api.BackupInstance) *api.BackupInstance) (*api.BackupInstance, kutil.VerbType, error) {
	return PatchBackupInstanceObject(c, cur, transform(cur.DeepCopy()))
}

func PatchBackupInstanceObject(c cs.StashV1alpha2Interface, cur, mod *api.BackupInstance) (*api.BackupInstance, kutil.VerbType, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(mod)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := jsonpatch.CreateMergePatch(curJson, modJson)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	glog.V(3).Infof("Patching BackupInstance %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.BackupInstances(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return out, kutil.VerbPatched, err
}

func TryUpdateBackupInstance(c cs.StashV1alpha2Interface, meta metav1.ObjectMeta, transform func(*api.BackupInstance) *api.BackupInstance) (result *api.BackupInstance, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.BackupInstances(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.BackupInstances(cur.Namespace).Update(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update BackupInstance %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update BackupInstance %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func SetBackupInstanceStats(c cs.StashV1alpha2Interface, backupInstnace *api.BackupInstance, backupStats restic_util.BackupStats, phase api.BackupInstancePhase) (*api.BackupInstance, error) {
	out, err := UpdateBackupInstanceStatus(c, backupInstnace, func(in *api.BackupInstanceStatus) *api.BackupInstanceStatus {
		in.Stats = backupStats
		in.Phase = phase
		return in
	}, apis.EnableStatusSubresource)
	return out, err
}

func UpdateBackupInstanceStatus(
	c cs.StashV1alpha2Interface,
	in *api.BackupInstance,
	transform func(*api.BackupInstanceStatus) *api.BackupInstanceStatus,
	useSubresource ...bool,
) (result *api.BackupInstance, err error) {
	if len(useSubresource) > 1 {
		return nil, errors.Errorf("invalid value passed for useSubresource: %v", useSubresource)
	}
	apply := func(x *api.BackupInstance) *api.BackupInstance {
		out := &api.BackupInstance{
			TypeMeta:   x.TypeMeta,
			ObjectMeta: x.ObjectMeta,
			Spec:       x.Spec,
			Status:     *transform(in.Status.DeepCopy()),
		}
		return out
	}

	if len(useSubresource) == 1 && useSubresource[0] {
		attempt := 0
		cur := in.DeepCopy()
		err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
			attempt++
			var e2 error
			result, e2 = c.BackupInstances(in.Namespace).UpdateStatus(apply(cur))
			if kerr.IsConflict(e2) {
				latest, e3 := c.BackupInstances(in.Namespace).Get(in.Name, metav1.GetOptions{})
				switch {
				case e3 == nil:
					cur = latest
					return false, nil
				case kutil.IsRequestRetryable(e3):
					return false, nil
				default:
					return false, e3
				}
			} else if err != nil && !kutil.IsRequestRetryable(e2) {
				return false, e2
			}
			return e2 == nil, nil
		})

		if err != nil {
			err = fmt.Errorf("failed to update status of BackupInstance %s/%s after %d attempts due to %v", in.Namespace, in.Name, attempt, err)
		}
		return
	}

	result, _, err = PatchBackupInstanceObject(c, in, apply(in))
	return
}
