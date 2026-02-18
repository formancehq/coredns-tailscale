package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	helmv2 "github.com/fluxcd/helm-controller/api/v2"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func newFluxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flux",
		Short: "FluxCD commands",
	}
	cmd.AddCommand(newFluxStatusCmd())
	return cmd
}

func newFluxStatusCmd() *cobra.Command {
	var server string
	var socks5 string
	var namespace string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show FluxCD reconciliation status",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := buildControllerRuntimeClient(server, socks5)
			if err != nil {
				return fmt.Errorf("building client: %w", err)
			}

			ctx := context.Background()
			listOpts := []ctrlclient.ListOption{}
			if namespace != "" {
				listOpts = append(listOpts, ctrlclient.InNamespace(namespace))
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "NAME\tKIND\tREADY\tMESSAGE\n")

			var kustomizations kustomizev1.KustomizationList
			if err := c.List(ctx, &kustomizations, listOpts...); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to list Kustomizations: %v\n", err)
			} else {
				for _, k := range kustomizations.Items {
					ready, msg := getReadyCondition(k.Status.Conditions)
					fmt.Fprintf(w, "%s/%s\tKustomization\t%s\t%s\n", k.Namespace, k.Name, ready, msg)
				}
			}

			var helmReleases helmv2.HelmReleaseList
			if err := c.List(ctx, &helmReleases, listOpts...); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to list HelmReleases: %v\n", err)
			} else {
				for _, h := range helmReleases.Items {
					ready, msg := getReadyCondition(h.Status.Conditions)
					fmt.Fprintf(w, "%s/%s\tHelmRelease\t%s\t%s\n", h.Namespace, h.Name, ready, msg)
				}
			}

			w.Flush()
			return nil
		},
	}

	cmd.Flags().StringVar(&server, "server", "", "Kubernetes API server URL (required)")
	cmd.Flags().StringVar(&socks5, "socks5", "", "SOCKS5 proxy address")
	cmd.Flags().StringVar(&namespace, "namespace", "flux-system", "Namespace to query")
	cmd.MarkFlagRequired("server")

	return cmd
}

func getReadyCondition(conditions []metav1.Condition) (string, string) {
	for _, c := range conditions {
		if c.Type == "Ready" {
			return string(c.Status), c.Message
		}
	}
	return "Unknown", ""
}
