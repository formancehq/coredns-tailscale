package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newPodsCmd() *cobra.Command {
	var server string
	var socks5 string
	var namespace string

	cmd := &cobra.Command{
		Use:   "pods",
		Short: "List pods",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientset, err := buildKubeClient(server, socks5)
			if err != nil {
				return fmt.Errorf("building kube client: %w", err)
			}

			pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
			if err != nil {
				return fmt.Errorf("listing pods: %w", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "NAMESPACE\tNAME\tSTATUS\n")
			for _, p := range pods.Items {
				fmt.Fprintf(w, "%s\t%s\t%s\n", p.Namespace, p.Name, string(p.Status.Phase))
			}
			w.Flush()
			return nil
		},
	}

	cmd.Flags().StringVar(&server, "server", "", "Kubernetes API server URL (required)")
	cmd.Flags().StringVar(&socks5, "socks5", "", "SOCKS5 proxy address")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (empty for all)")
	cmd.MarkFlagRequired("server")

	return cmd
}
