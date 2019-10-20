package main

import (
	"fmt"
	"log"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// DiffusionPodLabel Pod that checks the Label.
	DiffusionPodLabel = "client-type"
)

func main() {
	log.Print("Shared Informer app started")
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: onUpdate,
		AddFunc:    onAdd,
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

// onUpdate will check if Pod has been updated
func onAdd(obj interface{}) {
	// Cast the obj as node
	gpod := obj.(*corev1.Pod)
	_, ok := gpod.GetLabels()[DiffusionPodLabel]
	if ok {
		log.Printf("in Add!")
	}
}

// onUpdate will check if Pod has been updated
func onUpdate(obj interface{}, obj1 interface{}) {
	// Cast the obj as node
	gpod := obj.(*corev1.Pod)
	_, ok := gpod.GetLabels()[DiffusionPodLabel]
	if ok {
		log.Printf("in Update!")
	}
}
