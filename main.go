package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/load_balancers"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zones"
)

func main() {
	fmt.Println("[*] Fetching environment variables...")
	// Verify environment variables are set.
	accountID := os.Getenv("CF_ACCT_ID")
	if accountID == "" {
		fmt.Println("[-] ERROR - No CF_ACCT_ID provided.")
		os.Exit(1)
	}

	apiEmail := os.Getenv("CF_API_EMAIL")
	if apiEmail == "" {
		fmt.Println("[-] ERROR - No CF_API_EMAIL provided.")
		os.Exit(1)
	}

	apiKey := os.Getenv("CF_API_KEY")
	if apiKey == "" {
		fmt.Println("[-] ERROR - No CF_API_KEY provided.")
		os.Exit(1)
	}
	fmt.Println("[+] Environment variables fetched!")

	// Open Cloudflare client
	fmt.Println("[*] Opening Cloudflare Client...")
	client := cloudflare.NewClient(
		option.WithAPIKey(apiKey),
		option.WithAPIEmail(apiEmail),
	)

	// List zones in the account
	fmt.Println("[*] Listing Zones in Account ID " + accountID)
	zones, err1 := client.Zones.List(context.Background(), zones.ZoneListParams{
		Account: cloudflare.F(zones.ZoneListParamsAccount{ID: cloudflare.String(accountID)}),
	}, cloudflare.DefaultClientOptions()...)
	if err1 != nil {
		fmt.Println("[-] ERROR - Failed to list zones in Cloudflare Account.")
		fmt.Println(err1)
		os.Exit(1)
	}

	for _, zone := range zones.Result {
		fmt.Println("Fetching Load Balancers in Zone ID -- " + zone.ID)
		GetLBsInZone(client, zone.ID)
	}

}

func GetLBsInZone(client *cloudflare.Client, zoneID string) {
	lbs, err2 := client.LoadBalancers.List(context.Background(), load_balancers.LoadBalancerListParams{ZoneID: cloudflare.String(zoneID)}, cloudflare.DefaultClientOptions()...)
	if err2 != nil {
		fmt.Println("[-] ERROR - Failed to list load balancers in Cloudflare Zone.")
		fmt.Println(err2)
	}

	if lbs.Result == nil {
		fmt.Println("[/] No LBs in Cloudflare Zone.")
	}

	ItterateLBs(lbs.Result)

	// Could technically check sizes of pages before attempting to verify if there are additional pages.
	for {
		nextPage, err := lbs.GetNextPage()
		if err != nil {
			fmt.Println("[-] ERROR - Failed to list the next page of load balancers in Cloudflare Zone.")
			fmt.Println(err2)
		}

		if nextPage == nil {
			//fmt.Println("[/] No additional Load Balancer pages in Cloudflare Zone.")
			break
		} else {
			ItterateLBs(nextPage.Result)
		}
	}

}

func ItterateLBs(lbs []load_balancers.LoadBalancer) {
	for _, loadBalancer := range lbs {
		fmt.Printf("\tLoad Balancer ID: %s || Name: %s // Enabled: %v\r\n", loadBalancer.ID, loadBalancer.Name, loadBalancer.Enabled)
	}
}
