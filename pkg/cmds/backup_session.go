package cmds

import (
	"fmt"

	"github.com/appscode/go/log"

	//cs "github.com/appscode/stash/client/clientset/versioned"
	cs "github.com/appscode/stash/client/clientset/versioned/typed/stash/v1beta1"
	backup_session "github.com/appscode/stash/pkg/backup-session"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewBackupSession() *cobra.Command {
	var (
		masterURL      string
		kubeconfigPath string

		opt = backup_session.Options{
			//Namespace: meta.Namespace(),
		}
	)

	cmd := &cobra.Command{
		Use:               "backup_session",
		Short:             "backupSession CRD Create",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
			if err != nil {
				log.Fatalf("Could not get Kubernetes config: %s", err)
			}
			kubeClient := kubernetes.NewForConfigOrDie(config)
			stashv1beta1Client := cs.NewForConfigOrDie(config)

			ctrl := backup_session.New(kubeClient, stashv1beta1Client, opt)
			fmt.Println("hello...........01")
			err = ctrl.CreateBackupSessionCrd()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("hello...........02")

		},
	}

	cmd.Flags().StringVar(&masterURL, "master", "", "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&kubeconfigPath, "kubeconfig", "", "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVar(&opt.Name, "backupSessionName", "", "Set backupSessionName")
	cmd.Flags().StringVar(&opt.Namespace, "backupSessionNamespace", "", "Set backupSessionNamespace")

	return cmd
}
