package main

import (
	"fmt"
	"net"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	"golang.org/x/net/proxy"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func buildRestConfig(server, socks5Addr string) (*rest.Config, error) {
	config := &rest.Config{
		Host:            server,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	}
	if socks5Addr != "" {
		dialer, err := proxy.SOCKS5("tcp", socks5Addr, nil, &net.Dialer{Timeout: 30 * time.Second})
		if err != nil {
			return nil, fmt.Errorf("creating SOCKS5 dialer: %w", err)
		}
		config.Dial = dialer.(proxy.ContextDialer).DialContext
	}
	return config, nil
}

func buildKubeClient(server, socks5Addr string) (*kubernetes.Clientset, error) {
	config, err := buildRestConfig(server, socks5Addr)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func buildFluxScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kustomizev1.AddToScheme(scheme))
	utilruntime.Must(helmv2.AddToScheme(scheme))
	return scheme
}

func buildControllerRuntimeClient(server, socks5Addr string) (ctrlclient.Client, error) {
	config, err := buildRestConfig(server, socks5Addr)
	if err != nil {
		return nil, err
	}
	return ctrlclient.New(config, ctrlclient.Options{Scheme: buildFluxScheme()})
}
