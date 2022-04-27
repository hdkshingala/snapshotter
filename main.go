package main

import (
	"flag"
	"log"
	"path/filepath"
	"time"

	ssClient "github.com/hdkshingala/snapshotter/pkg/client/clientset/versioned"
	ssFactory "github.com/hdkshingala/snapshotter/pkg/client/informers/externalversions"
	ssController "github.com/hdkshingala/snapshotter/pkg/controller"
	snapshotclient "github.com/kubernetes-csi/external-snapshotter/client/v2/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("Building config from flags, %s", err.Error())
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Printf("Building config from InClusterConfig, %s", err.Error())
			return
		}
	}

	ssClientSet, err := ssClient.NewForConfig(config)
	if err != nil {
		log.Printf("Getting klient set, %s", err.Error())
	}

	vsClientSet, err := snapshotclient.NewForConfig(config)
	if err != nil {
		log.Printf("Getting client set, %s", err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Getting client set, %s", err.Error())
	}

	kInformerFactory := ssFactory.NewSharedInformerFactory(ssClientSet, 30*time.Second)

	k := ssController.NewController(clientSet, vsClientSet, ssClientSet, kInformerFactory.Hardik().V1alpha1().Snapshotters())

	ch := make(chan struct{})

	kInformerFactory.Start(ch)

	if err = k.Run(ch); err != nil {
		log.Printf("Error running controller, %s", err.Error())
	}

}
