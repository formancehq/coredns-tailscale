package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "internal-tool",
		Short: "Internal tool for Tailscale + Kubernetes + FluxCD diagnostics",
	}

	rootCmd.AddCommand(newKubeconfigCmd())
	rootCmd.AddCommand(newPodsCmd())
	rootCmd.AddCommand(newFluxCmd())
	rootCmd.AddCommand(newTailscaleCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
