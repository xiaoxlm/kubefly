package manager

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type Tunnel struct {
	protocol, proxyHost string
	cluster             string
	rate                int64
	timeout             time.Duration

	clientSet map[string]client.Client
}

func NewTunnel(cluster, protocol, proxyHost string, timeout int) *Tunnel {
	return &Tunnel{
		cluster:   cluster,
		protocol:  protocol,
		proxyHost: proxyHost,
		timeout:   time.Duration(timeout) * time.Second,
		clientSet: make(map[string]client.Client),
	}
}

func (tunnel *Tunnel) GetClient() (client.Client, error) {
	if tunnel.clientSet == nil {
		return nil, fmt.Errorf("client manager invalid")
	}

	if cli, ok := tunnel.getClient(); ok {
		return cli, nil
	}

	var cli client.Client
	{
		config, err := tunnel.GetProxyConfig()
		if err != nil {
			return nil, err
		}

		cli, err = client.New(config, client.Options{})
		if err != nil {
			return nil, err
		}
	}

	// cache
	tunnel.SetClient(cli)

	return cli, nil
}

func (tunnel *Tunnel) GetClientSet() (*kubernetes.Clientset, error) {
	config, err := tunnel.GetProxyConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func (tunnel *Tunnel) GetProxyConfig() (config *rest.Config, err error) {

	config, err = tunnel.remoteRESTConfig()
	if err != nil {
		return nil, err
	}

	config.Timeout = tunnel.timeout

	return config, nil
}

func (tunnel *Tunnel) getClient() (client.Client, bool) {
	cli, ok := tunnel.clientSet[tunnel.cluster]

	return cli, ok
}

func (tunnel *Tunnel) SetClient(cli client.Client) {
	tunnel.clientSet[tunnel.cluster] = cli
}

// todo
func (tunnel *Tunnel) remoteRESTConfig() (config *rest.Config, err error) {
	return
}
