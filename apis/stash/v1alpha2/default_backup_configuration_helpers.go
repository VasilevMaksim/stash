package v1alpha2

import (
	"hash/fnv"
	"strconv"

	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	hashutil "k8s.io/kubernetes/pkg/util/hash"
)

func (dbc DefaultBackupConfiguration) GetSpecHash() string {
	hash := fnv.New64a()
	hashutil.DeepHashObject(hash, dbc.Spec)
	return strconv.FormatUint(hash.Sum64(), 10)
}

func (dbc DefaultBackupConfiguration) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralDefaultBackupConfiguration,
		Singular:      ResourceSingularDefaultBackupConfiguration,
		Kind:          ResourceKindDefaultBackupConfiguration,
		ShortNames:    []string{"dbc"},
		Categories:    []string{"stash", "appscode", "backup"},
		ResourceScope: string(apiextensions.ClusterScoped),
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    SchemeGroupVersion.Version,
				Served:  true,
				Storage: true,
			},
		},
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "stash"},
		},
		SpecDefinitionName:    "github.com/appscode/stash/apis/stash/v1alpha2.DefaultBackupConfiguration",
		EnableValidation:      true,
		GetOpenAPIDefinitions: GetOpenAPIDefinitions,
		AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
			{
				Name:     "Backup-Procedure",
				Type:     "string",
				JSONPath: ".spec.backupProcedure",
			},
			{
				Name:     "Schedule",
				Type:     "string",
				JSONPath: ".spec.schedule",
			},
			{
				Name:     "Age",
				Type:     "date",
				JSONPath: ".metadata.creationTimestamp",
			},
		},
	})
}
