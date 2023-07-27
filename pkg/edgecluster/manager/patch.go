package manager

import (
	"encoding/json"
	"github.com/icza/dyno"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PatchFor(gvk schema.GroupVersionKind, live client.Object) client.Patch {
	if gvk.Group == corev1.GroupName && gvk.Kind == "Service" {
		return Merge(live)
	}

	// TODO handle more
	if gvk.Group == corev1.GroupName || gvk.Group == appsv1.GroupName {
		return client.StrategicMergeFrom(live)
	}

	return client.MergeFromWithOptions(live, client.MergeFromWithOptimisticLock{})
}

func Merge(live client.Object) client.Patch {
	return &MergePatch{
		live: live,
	}
}

// customer patch
type MergePatch struct {
	live client.Object

	applyObject client.Object
}

func (j *MergePatch) Type() types.PatchType {
	return types.MergePatchType
}

func (patch *MergePatch) Data(obj client.Object) ([]byte, error) {
	// init mergedObject
	patch.applyObject = obj

	mergedData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	merged := &unstructured.Unstructured{}
	mergedGVK := obj.GetObjectKind().GroupVersionKind()

	if _, _, err := unstructured.UnstructuredJSONScheme.Decode(mergedData, &mergedGVK, merged); err != nil {
		return nil, err
	}

	// filter
	patch.filter(patch.deleteObjectAnnotations)(merged.Object)

	return json.Marshal(merged.Object)
}

func (patch *MergePatch) filter(fns ...handleObjectUnstructured) func(objectUnstructured map[string]interface{}) {
	return func(objectUnstructured map[string]interface{}) {
		for _, f := range fns {
			f(objectUnstructured)
		}
	}
}

type handleObjectUnstructured func(objectUnstructured map[string]interface{})

func (patch *MergePatch) deleteObjectAnnotations(objectUnstructured map[string]interface{}) {
	liveAnnotation := patch.live.GetAnnotations()
	applyAnnotation := patch.applyObject.GetAnnotations()

	var deleteKeys []string
	for k := range liveAnnotation {
		if _, ok := applyAnnotation[k]; !ok {
			deleteKeys = append(deleteKeys, k)
		}
	}

	for _, k := range deleteKeys {
		_ = dyno.Set(objectUnstructured, nil, "metadata", "annotations", k)
	}
}
