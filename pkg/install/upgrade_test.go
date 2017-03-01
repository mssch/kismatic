package install

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/apprenda/kismatic/pkg/data"
)

type fakeUpgradeKubeClient struct {
	listPods                 func() (*data.PodList, error)
	getDaemonSet             func() (*data.DaemonSet, error)
	getReplicationController func() (*data.ReplicationController, error)
	getReplicaSet            func() (*data.ReplicaSet, error)
	getPersistentVolume      func(name string) (*data.PersistentVolume, error)
	getPersistentVolumeClaim func(name string) (*data.PersistentVolumeClaim, error)
	getStatefulSet           func() (*data.StatefulSet, error)
}

func (f fakeUpgradeKubeClient) ListPods() (*data.PodList, error) {
	if f.listPods != nil {
		return f.listPods()
	}
	return &data.PodList{}, nil
}

func (f fakeUpgradeKubeClient) GetDaemonSet(namespace, name string) (*data.DaemonSet, error) {
	if f.getDaemonSet != nil {
		return f.getDaemonSet()
	}
	return nil, errors.New("DS not found")
}

func (f fakeUpgradeKubeClient) GetReplicationController(namespace, name string) (*data.ReplicationController, error) {
	if f.getReplicationController != nil {
		return f.getReplicationController()
	}
	return nil, errors.New("RC not found")
}

func (f fakeUpgradeKubeClient) GetReplicaSet(namespace, name string) (*data.ReplicaSet, error) {
	if f.getReplicaSet != nil {
		return f.getReplicaSet()
	}
	return nil, errors.New("RS not found")
}

func (f fakeUpgradeKubeClient) GetPersistentVolume(name string) (*data.PersistentVolume, error) {
	if f.getPersistentVolume != nil {
		return f.getPersistentVolume(name)
	}
	return nil, errors.New("PV not found")
}

func (f fakeUpgradeKubeClient) GetPersistentVolumeClaim(namespace, name string) (*data.PersistentVolumeClaim, error) {
	if f.getPersistentVolumeClaim != nil {
		return f.getPersistentVolumeClaim(name)
	}
	return nil, errors.New("PVC not found")
}

func (f fakeUpgradeKubeClient) GetStatefulSet(namespace, name string) (*data.StatefulSet, error) {
	if f.getStatefulSet != nil {
		return f.getStatefulSet()
	}
	return nil, errors.New("StatefulSet not found")
}

func getSafePodWithCreatedByRef(t *testing.T, nodeName string, createdByKind string) data.Pod {
	createdByRef := data.SerializedReference{
		Reference: data.ObjectReference{
			Kind: createdByKind,
		},
	}
	b, err := json.Marshal(createdByRef)
	if err != nil {
		t.Fatalf("failed to marshal SerializedReference: %v", err)
	}
	return data.Pod{
		ObjectMeta: data.ObjectMeta{
			Name:        "foo",
			Namespace:   "foo",
			Annotations: map[string]string{kubeCreatedBy: string(b)},
		},
		Spec: data.PodSpec{
			NodeName: nodeName,
		},
	}
}

func TestDetectNodeUpgradeSafetyMasterCountUnsafe(t *testing.T) {
	plan := Plan{
		Master: MasterNodeGroup{
			ExpectedCount: 1,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Master.Nodes[0]
	k8sClient := fakeUpgradeKubeClient{}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if _, ok := errs[0].(masterNodeCountErr); !ok {
		t.Errorf("Expected masterNodeCountError, but got %v", errs[0])
	}
}

func TestDetectNodeUpgradeSafetyMasterLoadBalancingUnsafe(t *testing.T) {
	plan := Plan{
		Master: MasterNodeGroup{
			ExpectedCount:    2,
			LoadBalancedFQDN: "foo",
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
				{
					Host: "bar",
					IP:   "10.0.0.2",
				},
			},
		},
	}
	node := plan.Master.Nodes[0]
	k8sClient := fakeUpgradeKubeClient{}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if _, ok := errs[0].(masterNodeLoadBalancingErr); !ok {
		t.Errorf("Expected masterNodeCountError, but got %v", errs[0])
	}
}

func TestDetectNodeUpgradeSafetyMasterLoadBalancingSafe(t *testing.T) {
	plan := Plan{
		Master: MasterNodeGroup{
			ExpectedCount:    2,
			LoadBalancedFQDN: "someLoadBalancer",
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
				{
					Host: "bar",
					IP:   "10.0.0.2",
				},
			},
		},
	}
	node := plan.Master.Nodes[0]
	k8sClient := fakeUpgradeKubeClient{}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 0 {
		t.Errorf("did not expect an error, but got %d", len(errs))
	}
}

func TestDetectNodeUpgradeSafetyIngress(t *testing.T) {
	plan := Plan{
		Ingress: OptionalNodeGroup{
			ExpectedCount: 1,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Ingress.Nodes[0]
	k8sClient := fakeUpgradeKubeClient{}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if _, ok := errs[0].(ingressNotSupportedErr); !ok {
		t.Errorf("Expected ingressUnsupportedError, but got %v", errs[0])
	}
}

func TestDetectNodeUpgradeSafetyStorage(t *testing.T) {
	plan := Plan{
		Storage: OptionalNodeGroup{
			ExpectedCount: 1,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Storage.Nodes[0]
	k8sClient := fakeUpgradeKubeClient{}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if _, ok := errs[0].(storageNotSupportedErr); !ok {
		t.Errorf("Expected storageUnsupportedError, but got %v", errs[0])
	}
}

func TestDetectNodeUpgradeSafetyWorkerCountUnsafe(t *testing.T) {
	plan := Plan{
		Worker: NodeGroup{
			ExpectedCount: 1,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Worker.Nodes[0]
	k8sClient := fakeUpgradeKubeClient{}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if _, ok := errs[0].(workerNodeCountErr); !ok {
		t.Errorf("Expected storageUnsupportedError, but got %v", errs[0])
	}
}

func TestDetectNodeUpgradeSafetyWorkerPodListError(t *testing.T) {
	plan := Plan{
		Worker: NodeGroup{
			ExpectedCount: 2,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Worker.Nodes[0]
	k8sClient := fakeUpgradeKubeClient{listPods: func() (*data.PodList, error) { return nil, errors.New("some error") }}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if !strings.Contains(errs[0].Error(), "some error") {
		t.Errorf("unexpected error received: %v", errs[0])
	}
}

func TestDetectNodeUpgradeSafetyWorkerPodHostPathVol(t *testing.T) {
	plan := Plan{
		Worker: NodeGroup{
			ExpectedCount: 2,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Worker.Nodes[0]
	pod := getSafePodWithCreatedByRef(t, node.Host, "replicationcontroller")
	pod.Spec.Volumes = []data.Volume{
		{
			VolumeSource: data.VolumeSource{
				HostPath: &data.HostPathVolumeSource{},
			},
		},
	}
	k8sClient := fakeUpgradeKubeClient{
		listPods: func() (*data.PodList, error) {
			return &data.PodList{
				Items: []data.Pod{pod},
			}, nil
		},
		getReplicationController: func() (*data.ReplicationController, error) {
			return &data.ReplicationController{
				Status: data.ReplicationControllerStatus{
					Replicas: 2,
				},
			}, nil
		},
	}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if _, ok := errs[0].(podUnsafeVolumeErr); !ok {
		t.Errorf("expected podUnsafeVolumeError, but got %T", errs[0])
	}
}

func TestDetectNodeUpgradeSafetyWorkerPodEmptyDirVol(t *testing.T) {
	plan := Plan{
		Worker: NodeGroup{
			ExpectedCount: 2,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Worker.Nodes[0]

	// Setup a pod that was created by a daemonset
	pod := getSafePodWithCreatedByRef(t, node.Host, "ReplicationController")
	pod.Spec.Volumes = []data.Volume{
		{
			VolumeSource: data.VolumeSource{
				EmptyDir: &data.EmptyDirVolumeSource{},
			},
		},
	}
	k8sClient := fakeUpgradeKubeClient{
		listPods: func() (*data.PodList, error) {
			return &data.PodList{
				Items: []data.Pod{pod},
			}, nil
		},
		getReplicationController: func() (*data.ReplicationController, error) {
			return &data.ReplicationController{
				Status: data.ReplicationControllerStatus{
					Replicas: 2,
				},
			}, nil
		},
	}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if _, ok := errs[0].(podUnsafeVolumeErr); !ok {
		t.Errorf("expected podUnsafeVolumeError, but got %T", errs[0])
	}
}

func TestDetectNodeUpgradeSafetyWorkerPodHostPathPersistentVol(t *testing.T) {
	plan := Plan{
		Worker: NodeGroup{
			ExpectedCount: 2,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Worker.Nodes[0]
	pod := getSafePodWithCreatedByRef(t, node.Host, "replicationcontroller")
	// Setup a pod that is using a volume with a PersistentVolumeClaim
	// This PVC is bound to a persistent volume that is using HostPath
	pod.Spec.Volumes = []data.Volume{
		{
			VolumeSource: data.VolumeSource{
				PersistentVolumeClaim: &data.PersistentVolumeClaimVolumeSource{
					ClaimName: "theClaim",
				},
			},
		},
	}
	k8sClient := fakeUpgradeKubeClient{
		listPods: func() (*data.PodList, error) {
			return &data.PodList{
				Items: []data.Pod{pod},
			}, nil
		},
		getReplicationController: func() (*data.ReplicationController, error) {
			return &data.ReplicationController{
				Status: data.ReplicationControllerStatus{
					Replicas: 2,
				},
			}, nil
		},
		getPersistentVolumeClaim: func(name string) (*data.PersistentVolumeClaim, error) {
			if name == "theClaim" {
				return &data.PersistentVolumeClaim{
					Spec: data.PersistentVolumeClaimSpec{
						VolumeName: "theVolume",
					},
				}, nil
			}
			return nil, fmt.Errorf("PVC not found")
		},
		getPersistentVolume: func(name string) (*data.PersistentVolume, error) {
			if name == "theVolume" {
				spec := data.PersistentVolumeSpec{}
				spec.HostPath = &data.HostPathVolumeSource{}
				return &data.PersistentVolume{
					Spec: spec,
				}, nil
			}
			return nil, fmt.Errorf("PV not found")
		},
	}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if _, ok := errs[0].(podUnsafePersistentVolumeErr); !ok {
		t.Errorf("expected podUnsafePersistentVolumeError, but got %T", errs[0])
	}
}

func TestDetectNodeUpgradeSafetyWorkerSingleDaemon(t *testing.T) {
	plan := Plan{
		Worker: NodeGroup{
			ExpectedCount: 2,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Worker.Nodes[0]

	// Setup a pod that was created by a daemonset
	pod := getSafePodWithCreatedByRef(t, node.Host, "DaemonSet")

	// mock the k8s client to return the pod and an unsafe daemonset
	k8sClient := fakeUpgradeKubeClient{
		listPods: func() (*data.PodList, error) {
			return &data.PodList{
				Items: []data.Pod{pod},
			}, nil
		},
		getDaemonSet: func() (*data.DaemonSet, error) {
			return &data.DaemonSet{
				Status: data.DaemonSetStatus{
					DesiredNumberScheduled: 1,
				},
			}, nil
		},
	}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if _, ok := errs[0].(podUnsafeDaemonErr); !ok {
		t.Errorf("expected podSingleDaemonInstance, but got %T", errs[0])
	}
}

func TestDetectNodeUpgradeSafetyWorkerSafeDaemon(t *testing.T) {
	plan := Plan{
		Worker: NodeGroup{
			ExpectedCount: 2,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Worker.Nodes[0]

	// Setup a pod that was created by a daemonset
	pod := getSafePodWithCreatedByRef(t, node.Host, "DaemonSet")

	// mock the k8s client to return the pod and an safe daemonset
	k8sClient := fakeUpgradeKubeClient{
		listPods: func() (*data.PodList, error) {
			return &data.PodList{
				Items: []data.Pod{pod},
			}, nil
		},
		getDaemonSet: func() (*data.DaemonSet, error) {
			return &data.DaemonSet{
				Status: data.DaemonSetStatus{
					DesiredNumberScheduled: 2,
				},
			}, nil
		},
	}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 0 {
		t.Errorf("Did not expect errors, but got: %v", errs)
	}
}

func TestDetectNodeUpgradeSafetyLonePod(t *testing.T) {
	plan := Plan{
		Worker: NodeGroup{
			ExpectedCount: 2,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Worker.Nodes[0]

	// Setup a pod that is not managed by anything
	pod := getSafePodWithCreatedByRef(t, node.Host, "")
	delete(pod.Annotations, kubeCreatedBy)
	k8sClient := fakeUpgradeKubeClient{
		listPods: func() (*data.PodList, error) {
			return &data.PodList{
				Items: []data.Pod{pod},
			}, nil
		},
	}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if _, ok := errs[0].(unmanagedPodErr); !ok {
		t.Errorf("Expected a lonePodError, but got %v", errs[0])
	}
}

func TestDetectNodeUpgradeSafetyUnreplicatedController(t *testing.T) {
	plan := Plan{
		Worker: NodeGroup{
			ExpectedCount: 2,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Worker.Nodes[0]

	// Setup a pod that is managed by a replication controller with replicas = 1
	pod := getSafePodWithCreatedByRef(t, node.Host, "ReplicationController")
	fmt.Println(pod)
	k8sClient := fakeUpgradeKubeClient{
		listPods: func() (*data.PodList, error) {
			return &data.PodList{
				Items: []data.Pod{pod},
			}, nil
		},
		getReplicationController: func() (*data.ReplicationController, error) {
			return &data.ReplicationController{
				Status: data.ReplicationControllerStatus{
					Replicas: 1,
				},
			}, nil
		},
	}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if _, ok := errs[0].(unsafeReplicaCountErr); !ok {
		t.Errorf("expected unsafeReplicaCountErr, but got %T", errs[0])
	}
}

func TestDetectNodeUpgradeSafetyUnreplicatedReplicaSet(t *testing.T) {
	plan := Plan{
		Worker: NodeGroup{
			ExpectedCount: 2,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Worker.Nodes[0]

	// Setup a pod that is managed by a replication controller with replicas = 1
	pod := getSafePodWithCreatedByRef(t, node.Host, "ReplicaSet")
	k8sClient := fakeUpgradeKubeClient{
		listPods: func() (*data.PodList, error) {
			return &data.PodList{
				Items: []data.Pod{pod},
			}, nil
		},
		getReplicaSet: func() (*data.ReplicaSet, error) {
			return &data.ReplicaSet{
				Status: data.ReplicaSetStatus{
					Replicas: 1,
				},
			}, nil
		},
	}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if _, ok := errs[0].(unsafeReplicaCountErr); !ok {
		t.Errorf("expected unsafeReplicaCountErr, but got %T", errs[0])
	}
}

func TestDetectNodeUpgradeSafetyUnreplicatedStatefulSet(t *testing.T) {
	plan := Plan{
		Worker: NodeGroup{
			ExpectedCount: 2,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Worker.Nodes[0]

	// Setup a pod that is managed by a statefulset with replicas = 1
	pod := getSafePodWithCreatedByRef(t, node.Host, "StatefulSet")
	k8sClient := fakeUpgradeKubeClient{
		listPods: func() (*data.PodList, error) {
			return &data.PodList{
				Items: []data.Pod{pod},
			}, nil
		},
		getStatefulSet: func() (*data.StatefulSet, error) {
			return &data.StatefulSet{
				Status: data.StatefulSetStatus{
					Replicas: 1,
				},
			}, nil
		},
	}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if _, ok := errs[0].(unsafeReplicaCountErr); !ok {
		t.Errorf("expected unsafeReplicaCountErr, but got %T", errs[0])
	}
}

func TestDetectNodeUpgradeSafetyJobRunningOnNode(t *testing.T) {
	plan := Plan{
		Worker: NodeGroup{
			ExpectedCount: 2,
			Nodes: []Node{
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
				{
					Host: "foo",
					IP:   "10.0.0.1",
				},
			},
		},
	}
	node := plan.Worker.Nodes[0]

	// Setup a pod that is managed by a job
	pod := getSafePodWithCreatedByRef(t, node.Host, "Job")
	k8sClient := fakeUpgradeKubeClient{
		listPods: func() (*data.PodList, error) {
			return &data.PodList{
				Items: []data.Pod{pod},
			}, nil
		},
	}
	errs := DetectNodeUpgradeSafety(plan, node, k8sClient)
	if len(errs) != 1 {
		t.Errorf("Expected %d errors, but got %v", 1, errs)
	} else if _, ok := errs[0].(podRunningJobErr); !ok {
		t.Errorf("expected podRunningJobErr, but got %T", errs[0])
	}
}
