package discovery

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	servingv1alpha1 "github.com/knative/serving/pkg/client/clientset/versioned/typed/serving/v1alpha1"
	"github.com/golang/glog"
	"os"

	//istiov1alpha3 "github.com/knative/pkg/client/clientset/versioned/typed/istio/v1alpha3"

metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/api/core/v1"

	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
)

type KnativeFunction struct {
	ServiceName  string
	ServiceId    string
	HostName     string
	FunctionName string
	Revision 	string
	Namespace string
}

func DiscoverKnative(kubeClientConfig *restclient.Config) ([]*proto.EntityDTO, error) {

	fmt.Printf("[GO] hello, World\n")

	namespace := apiv1.NamespaceAll

	// knative serving client
	servingv1aplpha1Client, err := servingv1alpha1.NewForConfig(kubeClientConfig)
	if err != nil {
		fmt.Errorf("Error creating knative serving client: %++v", err)
		return nil, err
	}
	fmt.Printf("Got knative serving client\n")


	// Get knative services
	serviceList, err := servingv1aplpha1Client.Services(namespace).List(metav1.ListOptions{})
	if err != nil {
		fmt.Errorf("Error while getting services %++v", err)
		return nil, err
	}

	fmt.Printf("Got knative services %d\n", len(serviceList.Items))

	var functions []*KnativeFunction
	for _, svc := range serviceList.Items {

		fmt.Printf("*********** Service %s\n", svc.Name)
		//fmt.Printf("Service %++v\n", svc)
		function := & KnativeFunction{
			FunctionName: svc.Name,
			HostName: svc.Status.Domain,
			Revision: svc.Status.LatestReadyRevisionName,
			Namespace: svc.Namespace,
		}
		functions = append(functions, function)
	}

	// Build DTOs
	fmt.Printf("Building DTOs\n")
	knativeDtoBuilder := KnativeDTOBuilder{}
	var discoveryResult []*proto.EntityDTO
	for _, functionSvc := range functions {

		dtoBuilder, err := knativeDtoBuilder.buildFunctionDto(functionSvc)
		if err != nil {
			glog.Errorf("%s", err)
			fmt.Printf("Error while building entity : %v\n", err)
		}
		if dtoBuilder == nil {
			fmt.Printf("%v\n", err)
		}
		dto, err  := dtoBuilder.Create()
		if err != nil {
			fmt.Printf("builder error : %v\n", err)
		}
		fmt.Printf("Function DTO %++v\n", dto)
		discoveryResult = append(discoveryResult, dto)
	}
	fmt.Printf("DONE Building DTOs\n")
	return discoveryResult, nil
}

func createKubeConfigOrDie(kubeconfig string) *restclient.Config {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig) //TODO
	if err != nil {
		glog.Errorf("Fatal error: failed to get kubeconfig:  %s", err)
		os.Exit(1)
	}
	// This specifies the number and the max number of query per second to the api server.
	kubeConfig.QPS = 20.0
	kubeConfig.Burst = 30

	return kubeConfig
}

func createKubeClientOrDie(kubeConfig *restclient.Config) *kubernetes.Clientset {
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		glog.Errorf("Fatal error: failed to create kubeClient:%v", err)
		os.Exit(1)
	}

	return kubeClient
}

func GetNamespaces(kubeClient *kubernetes.Clientset) ([]*apiv1.Namespace, error) {

	namespaceList, err := kubeClient.Core().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list all namespaces in the cluster: %s", err)
	}
	namespaces := make([]*apiv1.Namespace, len(namespaceList.Items))
	for i := 0; i < len(namespaceList.Items); i++ {
		namespaces[i] = &namespaceList.Items[i]
		fmt.Printf("NAMESPACE %s\n", namespaces[i].Name)
	}
	return namespaces, nil
}
