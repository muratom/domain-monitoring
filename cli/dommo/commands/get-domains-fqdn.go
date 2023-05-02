package commands

import (
	"fmt"
	"net/http"
	"time"

	"github.com/muratom/domain-monitoring/api/rpc/v1/inspector"
	"github.com/spf13/cobra"
)

var getDomainsFQDNCmd = &cobra.Command{
	Use:   "get-domains-fqdn",
	Short: "Get all stored domains FQDN",
	Run: func(cmd *cobra.Command, args []string) {
		serverURI, err := cmd.Flags().GetString("server")
		if err != nil {
			fmt.Printf("faield to get 'server' flag")
			return
		}

		timeout, err := cmd.Flags().GetInt64("timeout")
		if err != nil {
			fmt.Printf("faield to get 'timeout' flag")
			return
		}

		req, err := inspector.NewGetAllDomainsRequest(serverURI)
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

		prettyJSON(resp.Body)
	},
}

func init() {
	getDomainsFQDNCmd.Flags().StringP("server", "s", "http://localhost:8000", "URI of a server to dial with")
	getDomainsFQDNCmd.Flags().Int64P("timeout", "t", 1000, "Server request timeout, ms")
}
