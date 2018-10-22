package discovery

import (
	servingv1alpha1 "github.com/knative/serving/pkg/client/clientset/versioned/typed/serving/v1alpha1"
	restclient "k8s.io/client-go/rest"
	"github.com/golang/glog"
)

func CreateKubeClientOrDie(kubeClientConfig *restclient.Config) (*servingv1alpha1.ServingV1alpha1Client, error){

	servingv1aplpha1Client, err := servingv1alpha1.NewForConfig(kubeClientConfig)
	if err != nil {
		glog.Errorf("Fatal error: failed to get knative client:  %s", err)
		return nil, err
	}
	return servingv1aplpha1Client, nil
}
