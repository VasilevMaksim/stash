package backup_session

import (
	"github.com/appscode/go/crypto/rand"
	api_v1beta1 "github.com/appscode/stash/apis/stash/v1beta1"

	//cs "github.com/appscode/stash/client/clientset/versioned"
	cs "github.com/appscode/stash/client/clientset/versioned/typed/stash/v1beta1"
	"github.com/appscode/stash/client/clientset/versioned/typed/stash/v1beta1/util"
	"github.com/tamalsaha/go-oneliners"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/apis/core"
)

type Controller struct {
	Options
	k8sClient          kubernetes.Interface
	stashv1beta1Client cs.StashV1beta1Interface
}

type Options struct {
	Name      string
	Namespace string
}

func New(k8sClient kubernetes.Interface, stashv1betaClient cs.StashV1beta1Interface, opt Options) *Controller {
	return &Controller{
		k8sClient:          k8sClient,
		stashv1beta1Client: stashv1betaClient,
		Options:            opt,
	}
}

func (c *Controller) CreateBackupSessionCrd() error {
	//timestamp := time.Now()
	//zoneName, zoneOffset := timestamp.Zone()
	//nameSuffix := strings.ToLower(fmt.Sprintf("-%d-%d-%d-%d-%d-%d-%s-%d",
	//	timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second(), zoneName, zoneOffset))

	backupSessionCrd := &api_v1beta1.BackupSession{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix(c.Name),
			Namespace: c.Namespace,
		},
		Spec: api_v1beta1.BackupSessionSpec{
			BackupConfiguration: &core.LocalObjectReference{
				Name: c.Name,
			},
		},
	}

	backupSession, _, err := util.CreateOrPatchBackupSession(c.stashv1beta1Client, backupSessionCrd.ObjectMeta, func(in *api_v1beta1.BackupSession) *api_v1beta1.BackupSession {

		if in.Spec.BackupConfiguration == nil {
			in.Spec.BackupConfiguration = &core.LocalObjectReference{}
		}
		in.Spec.BackupConfiguration.Name = backupSessionCrd.Spec.BackupConfiguration.Name

		return in

	})
	if err != nil {
		return err
	}
	oneliners.PrettyJson(backupSession)

	return nil

}
