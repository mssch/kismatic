package install

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apprenda/kismatic/pkg/data"
)

const kubeCreatedBy = "kubernetes.io/created-by"

type upgradeKubeInfoClient interface {
	data.PodLister
	data.DaemonSetGetter
	data.ReplicationControllerGetter
	data.ReplicaSetGetter
	data.PersistentVolumeClaimGetter
	data.PersistentVolumeGetter
	data.StatefulSetGetter
}

type masterNodeCountErr struct{}

func (e masterNodeCountErr) Error() string {
	return "This is the only master node in the cluster. " +
		"Upgrading it will make the cluster unavailable."
}

type masterNodeLoadBalancingErr struct{}

func (e masterNodeLoadBalancingErr) Error() string {
	return "This node is acting as the load balanced endpoint for the master nodes. " +
		"Upgrading it will make the cluster unavailable"
}

type ingressNotSupportedErr struct{}

func (e ingressNotSupportedErr) Error() string {
	return "Performing an online upgrade of ingress nodes is not supported."
}

type storageNotSupportedErr struct{}

func (e storageNotSupportedErr) Error() string {
	return "Performing an online upgrade of storage nodes is not supported."
}

type workerNodeCountErr struct{}

func (e workerNodeCountErr) Error() string {
	return "This is the only worker node in the cluster. " +
		"Upgrading it will make cluster features unavailable."
}

type podUnsafeVolumeErr struct {
	namespace string
	name      string
	volType   string
	volName   string
}

func (e podUnsafeVolumeErr) Error() string {
	return fmt.Sprintf(`Pod "%s/%s" is using %s volume %q, which is unsafe for upgrades.`, e.namespace, e.name, e.volType, e.volName)
}

type podUnsafePersistentVolumeErr struct {
	namespace string
	name      string
	volType   string
	volName   string
}

func (e podUnsafePersistentVolumeErr) Error() string {
	return fmt.Sprintf(`Pod "%s/%s" is using volume %q, which is backed by a %s PersistentVolume. `+
		`This kind of volume is unsafe for upgrades.`, e.namespace, e.name, e.volName, e.volType)
}

type podUnsafeDaemonErr struct {
	dsNamespace string
	dsName      string
}

func (e podUnsafeDaemonErr) Error() string {
	return fmt.Sprintf(`Pod managed by DaemonSet "%s/%s" is running on this node, and no other nodes `+
		"are capable of hosting this daemon. Upgrading it will make the daemon unavailable.", e.dsNamespace, e.dsName)
}

type unmanagedPodErr struct {
	namespace string
	name      string
}

func (e unmanagedPodErr) Error() string {
	return fmt.Sprintf(`The pod "%s/%s" is not being managed by a controller. `+
		"Upgrading this node might result in data or availability loss.", e.namespace, e.name)
}

type unsafeReplicaCountErr struct {
	kind      string
	namespace string
	name      string
}

func (e unsafeReplicaCountErr) Error() string {
	return fmt.Sprintf(`Pod managed by %s "%s/%s" is running on this node, `+
		"and the %s does not have a replica count greater than 1.", e.kind, e.namespace, e.name, e.kind)
}

type replicasOnSingleNodeErr struct {
	kind      string
	namespace string
	name      string
}

func (e replicasOnSingleNodeErr) Error() string {
	return fmt.Sprintf(`All the replicas that belong to the %s "%s/%s" are running on this node.`, e.kind, e.namespace, e.name)
}

type podRunningJobErr struct {
	namespace string
	name      string
}

func (e podRunningJobErr) Error() string {
	return fmt.Sprintf(`Pod that belongs to job "%s/%s" is running on this node.`, e.name, e.namespace)
}

// DetectNodeUpgradeSafety determines whether it's safe to upgrade a specific node
// listed in the plan file. If any condition that could result in data or availability
// loss is detected, the upgrade is deemed unsafe, and the conditions are returned as errors.
func DetectNodeUpgradeSafety(plan Plan, node Node, kubeClient upgradeKubeInfoClient) []error {
	errs := []error{}
	roles := plan.GetRolesForIP(node.IP)
	for _, role := range roles {
		switch role {
		case "master":
			if plan.Master.ExpectedCount < 2 {
				errs = append(errs, masterNodeCountErr{})
			}
			lbFQDN := plan.Master.LoadBalancedFQDN
			if lbFQDN == node.Host || lbFQDN == node.IP {
				errs = append(errs, masterNodeLoadBalancingErr{})
			}
		case "ingress":
			// we don't control load balancing of ingress nodes. therefore,
			// upgrading an ingress node is potentially unsafe
			errs = append(errs, ingressNotSupportedErr{})
		case "storage":
			// we could potentially detect safety of upgrading storage nodes by inspecting
			// the volumes on the node. for now, we are choosing not to support online upgrade of storage nodes
			errs = append(errs, storageNotSupportedErr{})
		case "worker":
			if plan.Worker.ExpectedCount < 2 {
				errs = append(errs, workerNodeCountErr{})
			}
			if workerErrs := detectWorkerNodeUpgradeSafety(node, kubeClient); workerErrs != nil {
				errs = append(errs, workerErrs...)
			}
		}
	}
	return errs
}

func detectWorkerNodeUpgradeSafety(node Node, kubeClient upgradeKubeInfoClient) []error {
	errs := []error{}
	podList, err := kubeClient.ListPods()
	if err != nil {
		errs = append(errs, fmt.Errorf("unable to determine node upgrade safety: %v", err))
		return errs
	}
	nodePods := []data.Pod{}
	for _, p := range podList.Items {
		if p.Spec.NodeName == node.Host {
			nodePods = append(nodePods, p)
		}
	}

	// Are there any pods using a hostPath, emptyDir volume OR a hostPath PersistentVolume?
	for _, p := range nodePods {
		for _, v := range p.Spec.Volumes {
			if v.VolumeSource.HostPath != nil {
				errs = append(errs, podUnsafeVolumeErr{namespace: p.Namespace, name: p.Name, volType: "HostPath", volName: v.Name})
			}
			if v.VolumeSource.EmptyDir != nil {
				errs = append(errs, podUnsafeVolumeErr{namespace: p.Namespace, name: p.Name, volType: "EmptyDir", volName: v.Name})
			}
			if v.VolumeSource.PersistentVolumeClaim != nil {
				claimRef := v.VolumeSource.PersistentVolumeClaim
				pvc, err := kubeClient.GetPersistentVolumeClaim(p.Namespace, claimRef.ClaimName)
				if err != nil {
					errs = append(errs, fmt.Errorf(`Failed to get PersistentVolumeClaim "%s/%s."`, p.Namespace, claimRef.ClaimName))
					continue
				}
				pvName := pvc.Spec.VolumeName
				pv, err := kubeClient.GetPersistentVolume(pvName)
				if err != nil {
					errs = append(errs, fmt.Errorf(`Failed to get PersistentVolume %q. This PV is being used by pod "%s/%s" on this node`, pvName, p.Namespace, p.Name))
					continue
				}
				if pv.Spec.HostPath != nil {
					errs = append(errs, podUnsafePersistentVolumeErr{namespace: p.Namespace, name: p.Name, volType: "HostPath", volName: v.Name})
				}
			}
		}
	}

	// Keep track of how many pods managed by replication controllers and replicasets
	// are running on this node. If all replicas are running on the node, we need to
	// return an error, as it would take the workload down.
	rcPods := map[string]int32{}
	rsPods := map[string]int32{}

	// 1. Are there any pods running on this node that are not managed by a controller?
	// 2. Are there any pods running on this node that are managed by a controller,
	//    and have replicas less than 2?
	// 3. Are there any daemonset managed pods running on this node? If so,
	//    verify that it is not the only one
	// 4. Are there any pods that belong to a job running on this node?
	for _, p := range nodePods {
		creator, ok := p.Annotations[kubeCreatedBy]
		if !ok {
			errs = append(errs, unmanagedPodErr{namespace: p.Namespace, name: p.Name})
			continue
		}
		var r data.SerializedReference
		err := json.Unmarshal([]byte(creator), &r)
		if err != nil {
			errs = append(errs, fmt.Errorf("Unable to determine the creator of pod %s/%s", p.Namespace, p.Name))
			continue
		}
		switch strings.ToLower(r.Reference.Kind) {
		default:
			errs = append(errs, fmt.Errorf("Unable to determine upgrade safety for a pod managed by a controller of type %q", r.Reference.Kind))
		case "daemonset":
			ds, err := kubeClient.GetDaemonSet(r.Reference.Namespace, r.Reference.Name)
			if err != nil {
				errs = append(errs, fmt.Errorf("DaemonSet pod is running, but failed to get information about DaemonSet %s/%s", r.Reference.Namespace, r.Reference.Name))
				continue
			}
			// Check if other nodes should be running this DS
			if ds.Status.DesiredNumberScheduled < 2 {
				errs = append(errs, podUnsafeDaemonErr{dsNamespace: r.Reference.Namespace, dsName: r.Reference.Name})
			}
		case "job":
			errs = append(errs, podRunningJobErr{namespace: r.Reference.Namespace, name: r.Reference.Name})
		case "replicationcontroller":
			rc, err := kubeClient.GetReplicationController(r.Reference.Namespace, r.Reference.Name)
			if err != nil {
				errs = append(errs, fmt.Errorf(`Failed to get information about replication controller "%s/%s"`, r.Reference.Namespace, r.Reference.Name))
			}
			if rc.Status.Replicas < 2 {
				errs = append(errs, unsafeReplicaCountErr{kind: r.Reference.Kind, namespace: r.Reference.Namespace, name: r.Reference.Name})
			}
			rcPods[r.Reference.Namespace+r.Reference.Name]++
			if rcPods[r.Reference.Namespace+r.Reference.Name] == rc.Status.Replicas {
				errs = append(errs, replicasOnSingleNodeErr{kind: r.Reference.Kind, namespace: r.Reference.Namespace, name: r.Reference.Name})
			}
		case "replicaset":
			rs, err := kubeClient.GetReplicaSet(r.Reference.Namespace, r.Reference.Name)
			if err != nil {
				errs = append(errs, fmt.Errorf(`Failed to get information about replica set "%s/%s"`, r.Reference.Namespace, r.Reference.Name))
			}
			if rs.Status.Replicas < 2 {
				errs = append(errs, unsafeReplicaCountErr{kind: r.Reference.Kind, namespace: r.Reference.Namespace, name: r.Reference.Name})
			}
			rsPods[r.Reference.Namespace+r.Reference.Name]++
			if rsPods[r.Reference.Namespace+r.Reference.Name] == rs.Status.Replicas {
				errs = append(errs, replicasOnSingleNodeErr{kind: r.Reference.Kind, namespace: r.Reference.Namespace, name: r.Reference.Name})
			}
		case "statefulset":
			sts, err := kubeClient.GetStatefulSet(r.Reference.Namespace, r.Reference.Name)
			if err != nil {
				errs = append(errs, fmt.Errorf(`Failed to get information about stateful set "%s/%s"`, r.Reference.Namespace, r.Reference.Name))
			}
			if sts.Status.Replicas < 2 {
				errs = append(errs, unsafeReplicaCountErr{kind: r.Reference.Kind, namespace: r.Reference.Namespace, name: r.Reference.Name})
			}
		}
	}

	return errs
}
