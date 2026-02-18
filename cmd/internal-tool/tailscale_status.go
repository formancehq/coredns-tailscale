package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"tailscale.com/client/local"
)

func newTailscaleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tailscale",
		Short: "Tailscale commands",
	}
	cmd.AddCommand(newTailscaleStatusCmd())
	return cmd
}

func newTailscaleStatusCmd() *cobra.Command {
	var socket string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show Tailscale peer status (requires access to tailscaled socket)",
		RunE: func(cmd *cobra.Command, args []string) error {
			lc := local.Client{
				Socket: socket,
			}

			status, err := lc.Status(context.Background())
			if err != nil {
				return fmt.Errorf("querying tailscale status: %w", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "HOSTNAME\tIP\tSTATUS\n")

			for _, peer := range status.Peer {
				state := "offline"
				if peer.Online {
					state = "online"
				}
				ip := ""
				if len(peer.TailscaleIPs) > 0 {
					ip = peer.TailscaleIPs[0].String()
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", peer.HostName, ip, state)
			}

			w.Flush()
			return nil
		},
	}

	cmd.Flags().StringVar(&socket, "socket", "", "Path to tailscaled socket (default: platform default)")

	return cmd
}
