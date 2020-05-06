# :rocket: Cloudflare DDNS

Dynamic DNS service based on Cloudflare! Access your home network remotely via a custom domain name without a static IP!

## :us: Origin

This script was written for the Raspberry Pi platform to enable low cost, simple self hosting to promote a more decentralized internet. On execution, the script fetches public IPv4 and IPv6 addresses and creates/updates DNS records for the subdomains in Cloudflare. Stale, duplicate DNS records are removed for housekeeping.

## :vertical_traffic_light: Getting Started

|    Environment Variable    | Default |     Example      |
| -------------------------- | ------- | ---------------- |
| `CLOUDFLARE_ACCOUNT_EMAIL` |         | user@example.com |
| `CLOUDFLARE_API_KEY`       |         |                  |
| `CLOUDFLARE_USE_PROXY`     | false   | false            |
| `CLOUDFLARE_ZONE_ID`       |         |                  |
| `CLOUDFLARE_SUBDOMAINS`    |         | ddns,home        |

Values explained:

```json
"CLOUDFLARE_API_KEY": "Your cloudflare API Key",
"CLOUDFLARE_ACCOUNT_EMAIL": "The email address you use to sign in to cloudflare",
"CLOUDFLARE_ZONE_ID": "The ID of the zone that will get the records. From your dashboard click into the zone. Under the overview tab, scroll down and the zone ID is listed in the right rail",
"CLOUDFLARE_SUBDOMAINS": "Array of subdomains you want to update the A & where applicable, AAAA records. IMPORTANT! Only write subdomain name. Do not include the base domain name. (e.g. foo or an empty string to update the base domain)",
"CLOUDFLARE_USE_PROXY": false (defaults to false. Make it true if you want CDN/SSL benefits from cloudflare. This usually disables SSH)
```

## :fax: Hosting multiple domains on the same IP?

You can save yourself some trouble when hosting multiple domains pointing to the same IP address (in the case of Traefik) by defining one A & AAAA record  'ddns.example.com' pointing to the IP of the server that will be updated by this DDNS script. For each subdomain, create a CNAME record pointing to 'ddns.example.com'. Now you don't have to manually modify the script config every time you add a new subdomain to your site!

## License

This Template is licensed under the GNU General Public License, version 3 (GPLv3) and is distributed free of charge.

## Author

Jacob McSwain

GitHub: https://github.com/USA-RedDragon

Website: https://jacob.mcswain.dev

## Original Author

Timothy Miller

GitHub: https://github.com/timothymiller 💡

Website: https://timknowsbest.com 💻

Donation: https://timknowsbest.com/donate 💸
