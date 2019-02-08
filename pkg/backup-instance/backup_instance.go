package backup_instance

import (
	"github.com/appscode/go/crypto/rand"
	api_v1alpha2 "github.com/appscode/stash/apis/stash/v1alpha2"
	cs "github.com/appscode/stash/client/clientset/versioned"
	"github.com/tamalsaha/go-oneliners"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
	//timestamp := time.Now()
	//zoneName, zoneOffset := timestamp.Zone()
	//nameSuffix := strings.ToLower(fmt.Sprintf("-%d-%d-%d-%d-%d-%d-%s-%d",
	//	timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second(), zoneName, zoneOffset))

	backupInstanceCrd := &api_v1alpha2.BackupInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix(c.Name),
			Namespace: c.Namespace,
		},
		Spec: api_v1alpha2.BackupInstanceSpec{
			TargetBackupConfiguration: c.Name,
		},
	}
	backupInstance, err := c.stashClient.StashV1alpha2().BackupInstances(c.Namespace).Create(backupInstanceCrd)
	if err != nil {
		return err
	}
	oneliners.PrettyJson(backupInstance)

	return nil

}
