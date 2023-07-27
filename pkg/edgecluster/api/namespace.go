package api

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetNamespace(ctx context.Context, clientset *kubernetes.Clientset, namespace string, opt metav1.GetOptions) (*corev1.Namespace, error) {
	return clientset.CoreV1().Namespaces().Get(ctx, namespace, opt)
}

func ListNamespace(ctx context.Context, clientset *kubernetes.Clientset, opt metav1.ListOptions) (*corev1.NamespaceList, error) {
	return clientset.CoreV1().Namespaces().List(ctx, opt)
}
