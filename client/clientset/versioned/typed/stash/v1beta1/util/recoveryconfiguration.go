package util

import (
	"fmt"
	"time"

	"github.com/appscode/kutil"
	"github.com/appscode/stash/apis"
	api "github.com/appscode/stash/apis/stash/v1beta1"
	cs "github.com/appscode/stash/client/clientset/versioned/typed/stash/v1beta1"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

func CreateOrPatchRecoveryConfiguration(c cs.StashV1beta1Interface, meta metav1.ObjectMeta, transform func(alert *api.RecoveryConfiguration) *api.RecoveryConfiguration) (*api.RecoveryConfiguration, kutil.VerbType, error) {
	cur, err := c.RecoveryConfigurations(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating RecoveryConfiguration %s/%s.", meta.Namespace, meta.Name)
		out, err := c.RecoveryConfigurations(meta.Namespace).Create(transform(&api.RecoveryConfiguration{
			TypeMeta: metav1.TypeMeta{
				Kind:       "RecoveryConfiguration",
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchRecoveryConfiguration(c, cur, transform)
}

func PatchRecoveryConfiguration(c cs.StashV1beta1Interface, cur *api.RecoveryConfiguration, transform func(*api.RecoveryConfiguration) *api.RecoveryConfiguration) (*api.RecoveryConfiguration, kutil.VerbType, error) {
	return PatchRecoveryConfigurationObject(c, cur, transform(cur.DeepCopy()))
}

func PatchRecoveryConfigurationObject(c cs.StashV1beta1Interface, cur, mod *api.RecoveryConfiguration) (*api.RecoveryConfiguration, kutil.VerbType, error) {
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
	glog.V(3).Infof("Patching RecoveryConfiguration %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.RecoveryConfigurations(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return out, kutil.VerbPatched, err
}

func TryUpdateRecoveryConfiguration(c cs.StashV1beta1Interface, meta metav1.ObjectMeta, transform func(*api.RecoveryConfiguration) *api.RecoveryConfiguration) (result *api.RecoveryConfiguration, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.RecoveryConfigurations(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.RecoveryConfigurations(cur.Namespace).Update(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update RecoveryConfiguration %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update RecoveryConfiguration %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func SetRecoveryConfigurationStats(c cs.StashV1beta1Interface, recovery *api.RecoveryConfiguration, path string, d time.Duration, phase api.RecoveryPhase) (*api.RecoveryConfiguration, error) {
	out, err := UpdateRecoveryConfigurationStatus(c, recovery, func(in *api.RecoveryConfigurationStatus) *api.RecoveryConfigurationStatus {
		found := false
		for _, stats := range in.Stats {
			if stats.Path == path {
				found = true
				stats.Duration = d.String()
				stats.Phase = phase
			}
		}
		if !found {
			recovery.Status.Stats = append(recovery.Status.Stats, api.RestoreStats{
				Path:     path,
				Duration: d.String(),
				Phase:    phase,
			})
		}
		return in
	}, apis.EnableStatusSubresource)
	return out, err
}

func UpdateRecoveryConfigurationStatus(
	c cs.StashV1beta1Interface,
	in *api.RecoveryConfiguration,
	transform func(*api.RecoveryConfigurationStatus) *api.RecoveryConfigurationStatus,
	useSubresource ...bool,
) (result *api.RecoveryConfiguration, err error) {
	if len(useSubresource) > 1 {
		return nil, errors.Errorf("invalid value passed for useSubresource: %v", useSubresource)
	}
	apply := func(x *api.RecoveryConfiguration) *api.RecoveryConfiguration {
		out := &api.RecoveryConfiguration{
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
			result, e2 = c.RecoveryConfigurations(in.Namespace).UpdateStatus(apply(cur))
			if kerr.IsConflict(e2) {
				latest, e3 := c.RecoveryConfigurations(in.Namespace).Get(in.Name, metav1.GetOptions{})
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
			err = fmt.Errorf("failed to update status of RecoveryConfiguration %s/%s after %d attempts due to %v", in.Namespace, in.Name, attempt, err)
		}
		return
	}

	result, _, err = PatchRecoveryConfigurationObject(c, in, apply(in))
	return
}
