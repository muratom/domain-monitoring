package commands

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/muratom/domain-monitoring/api/rpc/v1/inspector"
	"github.com/muratom/domain-monitoring/api/rpc/v1/inspector/models"
	"github.com/spf13/cobra"
)

var addDomainCmd = &cobra.Command{
	Use:   "add-domain",
	Short: "Add new domain by FQDN",
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

		params := &models.AddDomainParams{
			Fqdn: string(fqdn),
		}
		req, err := inspector.NewAddDomainRequest(serverURI, params)
		if err != nil {
			fmt.Printf("failed to form a request: %v", err)
			return
		}

		client := http.Client{
			Timeout: time.Millisecond * time.Duration(timeout),
		}
		resp, err := client.Do(req)
		if err != nil {
			if os.IsTimeout(err) {
				fmt.Println("request timeout happened")
			} else {
				fmt.Printf("error doing a request to server (%v): %v", serverURI, err)
			}
			return
		}

		prettyJSON(resp.Body)
	},
}

func init() {
	addDomainCmd.Flags().String("fqdn", "", "Requested fully-qualified domain name")
	addDomainCmd.Flags().StringP("server", "s", "http://localhost:8000", "URI of a server to dial with")
	addDomainCmd.Flags().Int64P("timeout", "t", 1000, "Server request timeout, ms")
}
