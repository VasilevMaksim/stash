package backup_instance

import (
	"fmt"
	api_v1alpha2 "github.com/appscode/stash/apis/stash/v1alpha2"
	cs "github.com/appscode/stash/client/clientset/versioned"
	"github.com/tamalsaha/go-oneliners"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"math/rand"
	"strconv"
)

type Controller struct {
	Options
	stashClient cs.Interface
	k8sClient   kubernetes.Interface
}

type Options struct {
	Name      string
	Namespace string
}

func New(k8sClient kubernetes.Interface, stashClient cs.Interface, opt Options) *Controller {
	return &Controller{
		k8sClient:   k8sClient,
		stashClient: stashClient,
		Options:     opt,
	}
}

func (c *Controller) CreateBackupInstanceCrd() error {
	backupInstanceCrd := &api_v1alpha2.BackupInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name + "-" + strconv.Itoa(rand.Int()),
			Namespace: c.Namespace,
		},
		Spec: api_v1alpha2.BackupInstanceSpec{
			TargetBackupConfiguration: c.Name,
		},
	}
	fmt.Println("backupInastance:::::::::::::::::")
	backupInstance, err := c.stashClient.StashV1alpha2().BackupInstances(c.Namespace).Create(backupInstanceCrd)
	if err != nil {
		return err
	}
	oneliners.PrettyJson(backupInstance)

	return nil

}
