package api

import (
	"context"
	"github.com/xiaoxlm/kubefly/pkg/edgecluster/manager"
	"github.com/xiaoxlm/kubefly/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"testing"
)

var (
	ctx       = context.Background()
	clientset *kubernetes.Clientset
)

func init() {
	var (
		err           error
		clientManager = manager.NewTunnel("" /*cluster name*/, "" /*http or https*/, "" /*tunnel host*/, 60)
	)
	clientset, err = clientManager.GetClientSet()
	if err != nil {
		panic(err)
	}
}

func TestGetNamespace(t *testing.T) {
	type args struct {
		ctx       context.Context
		clientset *kubernetes.Clientset
		namespace string
		opt       metav1.GetOptions
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "#devops",
			args: args{
				ctx:       ctx,
				clientset: clientset,
				namespace: "devops",
				opt:       metav1.GetOptions{},
			},
			want:    "devops",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNamespace(tt.args.ctx, tt.args.clientset, tt.args.namespace, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.GetObjectMeta().GetName() != tt.want {
				t.Errorf("GetNamespace() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListNamespace(t *testing.T) {
	got, err := ListNamespace(ctx, clientset, metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}

	util.LogJSON(got)
}
