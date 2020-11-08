package main

import (
	"fmt"
	"os"
	"time"

	"io/ioutil"
	"net/http"

	"k8s.io/klog/v2"

	"github.com/cloudflare/cloudflare-go"
)

func fetchIP(v6 bool) string {
	publicIPUrl := "https://api.ipify.org?format=text"

	if v6 {
		publicIPUrl = "https://api6.ipify.org?format=text"
	}

	klog.Info("Fetching IP address from external API")
	resp, err := http.Get(publicIPUrl)
	if err != nil {
		klog.Exit(err)
	}
	defer resp.Body.Close()

	ipBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Exit(err)
	}
	return string(ipBytes)
}

func upsertIP(cfAPI *cloudflare.API, zoneID string, newRecord *cloudflare.DNSRecord, create bool) error {
	if create {
		rr, err := cfAPI.CreateDNSRecord(zoneID, *newRecord)
		if err != nil {
			return err
		}
		if !rr.Response.Success {
			klog.Exitf("Failed to create record: %v", rr.Response)
		}
	} else {
		err := cfAPI.UpdateDNSRecord(zoneID, newRecord.ID, *newRecord)
		if err != nil {
			return err
		}
	}
	return nil
}

func run(cfAPI *cloudflare.API, zoneID string, cfRecord string, ipv6 bool) {
	for {
		needCreateRecord := false
		recordType := "A"
		if ipv6 {
			recordType = "AAAA"
		}
		currentRecords, err := cfAPI.DNSRecords(zoneID, cloudflare.DNSRecord{Name: cfRecord, Type: recordType})
		if err != nil {
			klog.Exit(err)
		}
		if len(currentRecords) > 1 {
			klog.Exit("Multiple type A records returned with the same name.")
		} else if len(currentRecords) == 0 {
			klog.Info("No existing record found.")
			needCreateRecord = true
		}

		ip := fetchIP(ipv6)
		klog.Infof("Public IP is: %s", ip)

		if !needCreateRecord {
			klog.Infof("CF Record IP is: %s", currentRecords[0].Content)
		}

		if needCreateRecord || ip != currentRecords[0].Content {
			klog.Info("Updating record due to IP mismatch")

			newRecord := cloudflare.DNSRecord{
				Name:    cfRecord,
				Content: ip,
				Type:    recordType,
				TTL:     120,
			}
			if !needCreateRecord {
				newRecord.ID = currentRecords[0].ID
			}

			err = upsertIP(cfAPI, zoneID, &newRecord, needCreateRecord)
			if err != nil {
				klog.Exit(err)
			}

			klog.Info("Updated DNS record")
		}
		time.Sleep(30 * time.Second)
	}
}

func main() {
	cfAPIToken := os.Getenv("CF_API_TOKEN")
	cfZone := os.Getenv("CF_ZONE")
	cfRecord := fmt.Sprintf("%s.%s", os.Getenv("CF_RECORD"), cfZone)
	ipv6Str := os.Getenv("IPV6")
	ipv6 := false

	cfAPI, err := cloudflare.NewWithAPIToken(cfAPIToken)
	if err != nil {
		klog.Exit(err)
	}

	zoneID, err := cfAPI.ZoneIDByName(cfZone)
	if err != nil {
		klog.Exit(err)
	}

	if ipv6Str != "" {
		ipv6 = true
	}

	go run(cfAPI, zoneID, cfRecord, ipv6)
	if ipv6 {
		// Counter-intuitively, this !ipv6 ensures that IPv4 addresses are updated
		// In addition to the v6 addresses
		go run(cfAPI, zoneID, cfRecord, !ipv6)
	}

	for {
		c := make(chan int)
		<-c
	}

}
