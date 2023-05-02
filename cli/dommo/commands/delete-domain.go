package commands

import (
	"fmt"
	"net/http"
	"time"

	"github.com/muratom/domain-monitoring/api/rpc/v1/inspector"
	"github.com/muratom/domain-monitoring/api/rpc/v1/inspector/models"
	"github.com/spf13/cobra"
)

var deleteDomainCmd = &cobra.Command{
	Use:   "delete-domain",
	Short: "Delete a domain from storage by FQDN",
	Run: func(cmd *cobra.Command, args []string) {
		fqdn, err := cmd.Flags().GetString("fqdn")
		if err != nil {
			fmt.Printf("failed to get 'fqdn' flag")
			return
		}
		utfFQDN := []rune(fqdn)
		if len(utfFQDN) == 0 || utfFQDN[len(utfFQDN)-1] != '.' {
			fmt.Println("invalid FQDN: it must be not empty and ends with a dot '.'")
			return
		}

		serverURI, err := cmd.Flags().GetString("server")
		if err != nil {
			fmt.Printf("failed to get 'server' flag")
			return
		}

		timeout, err := cmd.Flags().GetInt64("timeout")
		if err != nil {
			fmt.Printf("failed to get 'timeout' flag")
			return
		}

		params := &models.DeleteDomainParams{
			Fqdn: string(fqdn),
		}
		req, err := inspector.NewDeleteDomainRequest(serverURI, params)
		if err != nil {
			fmt.Printf("failed to form a request: %v", err)
			return
		}

		client := http.Client{
			Timeout: time.Millisecond * time.Duration(timeout),
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("error doing a request to server (%v): %v", serverURI, err)
			return
		}

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("OK: domain %v deleted\n", fqdn)
			return
		}

		fmt.Printf("error deleting domain: %v\n", resp.StatusCode)
	},
}

func init() {
	deleteDomainCmd.Flags().String("fqdn", "", "Requested fully-qualified domain name")
	deleteDomainCmd.Flags().StringP("server", "s", "http://localhost:8000", "URI of a server to dial with")
	deleteDomainCmd.Flags().Int64P("timeout", "t", 1000, "Server request timeout, ms")
}
