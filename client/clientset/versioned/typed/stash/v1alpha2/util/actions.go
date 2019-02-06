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

func CreateOrPatchAction(c cs.StashV1alpha2Interface, meta metav1.ObjectMeta, transform func(alert *api.Action) *api.Action) (*api.Action, kutil.VerbType, error) {
	cur, err := c.Actions().Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating Action %s/%s.", meta.Namespace, meta.Name)
		out, err := c.Actions().Create(transform(&api.Action{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Action",
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchAction(c, cur, transform)
}

func PatchAction(c cs.StashV1alpha2Interface, cur *api.Action, transform func(*api.Action) *api.Action) (*api.Action, kutil.VerbType, error) {
	return PatchActionObject(c, cur, transform(cur.DeepCopy()))
}

func PatchActionObject(c cs.StashV1alpha2Interface, cur, mod *api.Action) (*api.Action, kutil.VerbType, error) {
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
	glog.V(3).Infof("Patching Action %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.Actions().Patch(cur.Name, types.MergePatchType, patch)
	return out, kutil.VerbPatched, err
}

func TryUpdateAction(c cs.StashV1alpha2Interface, meta metav1.ObjectMeta, transform func(*api.Action) *api.Action) (result *api.Action, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Actions().Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.Actions().Update(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update Action %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update Action %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}
