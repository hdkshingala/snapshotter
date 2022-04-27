package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hdkshingala/snapshotter/pkg/apis/hardik.dev/v1alpha1"
	snapshotv1beta1 "github.com/kubernetes-csi/external-snapshotter/client/v2/apis/volumesnapshot/v1beta1"

	"github.com/hdkshingala/snapshotter/pkg/client/clientset/versioned/scheme"

	snapshotclient "github.com/kubernetes-csi/external-snapshotter/client/v2/clientset/versioned"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	ssClientset "github.com/hdkshingala/snapshotter/pkg/client/clientset/versioned"
	ssInformer "github.com/hdkshingala/snapshotter/pkg/client/informers/externalversions/hardik.dev/v1alpha1"
	ssLister "github.com/hdkshingala/snapshotter/pkg/client/listers/hardik.dev/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	client          kubernetes.Interface
	vsClient        snapshotclient.Interface
	ssClient        ssClientset.Interface
	ssClusterSynced cache.InformerSynced
	ssLister        ssLister.SnapshotterLister
	wq              workqueue.RateLimitingInterface
	recorder        record.EventRecorder
}

func NewController(client kubernetes.Interface, vsClient snapshotclient.Interface, ssClient ssClientset.Interface, informer ssInformer.SnapshotterInformer) *Controller {
	runtime.Must(scheme.AddToScheme(scheme.Scheme))

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{
		Interface: client.CoreV1().Events(""),
	})

	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{
		Component: "Snapshotter",
	})

	cont := &Controller{
		client:          client,
		vsClient:        vsClient,
		ssClient:        ssClient,
		ssClusterSynced: informer.Informer().HasSynced,
		ssLister:        informer.Lister(),
		wq:              workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Snapshotter"),
		recorder:        recorder,
	}

	informer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: cont.handleAdd,
		},
	)

	return cont
}

func (cont *Controller) handleAdd(obj interface{}) {
	cont.wq.Add(obj.(*v1alpha1.Snapshotter))
}

func (cont *Controller) Run(ch chan struct{}) error {
	if bool := cache.WaitForCacheSync(ch, cont.ssClusterSynced); !bool {
		log.Println("cache was not synced.")
	}

	go wait.Until(cont.worker, time.Second, ch)
	<-ch
	return nil
}

func (cont *Controller) worker() {
	for cont.process() {
	}
}

func (cont *Controller) process() bool {
	item, shutDown := cont.wq.Get()
	if shutDown {
		log.Println("Cache is closed")
		return false
	}
	defer cont.wq.Forget(item)

	ss, ok := item.(*v1alpha1.Snapshotter)
	if !ok {
		log.Printf("Error received while converting the object '%+v' to Snapshotter type.", item)
		return false
	}

	log.Printf("Snapshotter spec is '%+v'", ss.Spec)

	err := cont.createVolumeSnapshot(ss)
	if err != nil {
		log.Printf("Error received which creating cluster, %s", err.Error())
		cont.recorder.Event(
			ss, corev1.EventTypeNormal, "SnapShotCreationFailed", fmt.Sprintf("Snapshot with name '%s-snapshot' was failed due to '%s'.", ss.Spec.ClaimName, err.Error()))
		return false
	}

	cont.recorder.Event(ss, corev1.EventTypeNormal, "SnapShotCreationCompleted", fmt.Sprintf("Snapshot with name '%s-snapshot' was created.", ss.Spec.ClaimName))

	return true
}

func (cont *Controller) createVolumeSnapshot(ss *v1alpha1.Snapshotter) error {
	snapshot := snapshotv1beta1.VolumeSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ss.Spec.ClaimName + "-snapshot",
			Namespace: ss.Namespace,
		},
		Spec: snapshotv1beta1.VolumeSnapshotSpec{
			VolumeSnapshotClassName: &ss.Spec.ClassName,
			Source: snapshotv1beta1.VolumeSnapshotSource{
				PersistentVolumeClaimName: &ss.Spec.ClaimName,
			},
		},
	}
	_, err := cont.vsClient.SnapshotV1beta1().VolumeSnapshots(ss.Namespace).Create(context.Background(), &snapshot, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
