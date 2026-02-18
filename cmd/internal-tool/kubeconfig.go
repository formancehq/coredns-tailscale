package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newKubeconfigCmd() *cobra.Command {
	var server string
	var socks5 string

	cmd := &cobra.Command{
		Use:   "kubeconfig",
		Short: "Verify Kubernetes API connectivity via Tailscale proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientset, err := buildKubeClient(server, socks5)
			if err != nil {
				return fmt.Errorf("building kube client: %w", err)
			}

			version, err := clientset.Discovery().ServerVersion()
			if err != nil {
				return fmt.Errorf("querying server version: %w", err)
			}

			fmt.Printf("Connected to Kubernetes API\n")
			fmt.Printf("  Server: %s\n", server)
			fmt.Printf("  Version: %s\n", version.GitVersion)
			fmt.Printf("  Platform: %s\n", version.Platform)
			return nil
		},
	}

	cmd.Flags().StringVar(&server, "server", "", "Kubernetes API server URL (required)")
	cmd.Flags().StringVar(&socks5, "socks5", "", "SOCKS5 proxy address (e.g., ts-proxy:1055)")
	cmd.MarkFlagRequired("server")

	return cmd
}
