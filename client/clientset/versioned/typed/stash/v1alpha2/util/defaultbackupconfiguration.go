package util

import (
	"fmt"

	"github.com/appscode/kutil"
	api "github.com/appscode/stash/apis/stash/v1alpha2"
	cs "github.com/appscode/stash/client/clientset/versioned/typed/stash/v1alpha2"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/golang/glog"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

func CreateOrPatchDefaultBackupConfiguration(c cs.StashV1alpha2Interface, meta metav1.ObjectMeta, transform func(alert *api.DefaultBackupConfiguration) *api.DefaultBackupConfiguration) (*api.DefaultBackupConfiguration, kutil.VerbType, error) {
	cur, err := c.DefaultBackupConfigurations().Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating DefaultBackupConfiguration %s/%s.", meta.Namespace, meta.Name)
		out, err := c.DefaultBackupConfigurations().Create(transform(&api.DefaultBackupConfiguration{
			TypeMeta: metav1.TypeMeta{
				Kind:       "DefaultBackupConfiguration",
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchDefaultBackupConfiguration(c, cur, transform)
}

func PatchDefaultBackupConfiguration(c cs.StashV1alpha2Interface, cur *api.DefaultBackupConfiguration, transform func(*api.DefaultBackupConfiguration) *api.DefaultBackupConfiguration) (*api.DefaultBackupConfiguration, kutil.VerbType, error) {
	return PatchDefaultBackupConfigurationObject(c, cur, transform(cur.DeepCopy()))
}

func PatchDefaultBackupConfigurationObject(c cs.StashV1alpha2Interface, cur, mod *api.DefaultBackupConfiguration) (*api.DefaultBackupConfiguration, kutil.VerbType, error) {
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
	glog.V(3).Infof("Patching DefaultBackupConfiguration %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.DefaultBackupConfigurations().Patch(cur.Name, types.MergePatchType, patch)
	return out, kutil.VerbPatched, err
}

func TryUpdateDefaultBackupConfiguration(c cs.StashV1alpha2Interface, meta metav1.ObjectMeta, transform func(*api.DefaultBackupConfiguration) *api.DefaultBackupConfiguration) (result *api.DefaultBackupConfiguration, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.DefaultBackupConfigurations().Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.DefaultBackupConfigurations().Update(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update DefaultBackupConfiguration %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update DefaultBackupConfiguration %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}
