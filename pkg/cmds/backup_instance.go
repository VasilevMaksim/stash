package cmds

import (
	"github.com/appscode/go/log"
	"github.com/appscode/kutil/meta"
	cs "github.com/appscode/stash/client/clientset/versioned"
	backup_instance "github.com/appscode/stash/pkg/backup-instance"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewBackupInstance() *cobra.Command {
	var (
		masterURL      string
		kubeconfigPath string

		opt = backup_instance.Options{
			Namespace: meta.Namespace(),
		}
	)

	cmd := &cobra.Command{
		Use:               "backup_instance",
		Short:             "backupInastance CRD Create",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
			if err != nil {
				log.Fatalf("Could not get Kubernetes config: %s", err)
			}
			kubeClient := kubernetes.NewForConfigOrDie(config)
			stashClient := cs.NewForConfigOrDie(config)

			ctrl := backup_instance.New(kubeClient, stashClient, opt)
			err = ctrl.CreateBackupInstanceCrd()
			if err != nil {
				log.Fatal(err)
			}

		},
	}

	cmd.Flags().StringVar(&masterURL, "master", "", "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&kubeconfigPath, "kubeconfig", "", "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVar(&opt.Name, "backupInstanceName", "", "Set backupInstanceName")

	return cmd
}
