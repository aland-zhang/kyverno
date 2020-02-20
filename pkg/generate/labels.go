package generate

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func manageLabels(unstr *unstructured.Unstructured, triggerResource unstructured.Unstructured) {
	// add managedBY label if not defined
	labels := unstr.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}

	// handle managedBy label
	managedBy(labels)
	// handle generatedBy label
	generatedBy(labels, triggerResource)

	// update the labels
	unstr.SetLabels(labels)
}

func managedBy(labels map[string]string) {
	// ManagedBy label
	key := "app.kubernetes.io/managed-by"
	value := "kyverno"
	val, ok := labels[key]
	if ok {
		if val != value {
			glog.Infof("resource managed by %s, kyverno wont over-ride the label", val)
			return
		}
	}
	if !ok {
		// add label
		labels[key] = value
	}
}

func generatedBy(labels map[string]string, triggerResource unstructured.Unstructured) {
	key := "kyverno.io/generated-by"
	value := fmt.Sprintf("%s-%s-%s", triggerResource.GetKind(), triggerResource.GetNamespace(), triggerResource.GetName())
	val, ok := labels[key]
	if ok {
		if val != value {
			glog.Infof("resource generated by %s, kyverno wont over-ride the label", val)
			return
		}
	}
	if !ok {
		// add label
		labels[key] = value
	}
}