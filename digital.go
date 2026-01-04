package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)



type InputType string

const (
	InputIP       InputType = "ip"
	InputUsername InputType = "username"
)

type ValidationResult struct {
	Type  InputType
	Value string
}

func ValidatePublicInput(raw string) (*ValidationResult, error) {
	clean := strings.TrimSpace(raw)

	if clean == "" {
		return nil, errors.New("input is empty")
	}

	if ip := net.ParseIP(clean); ip != nil {
		return &ValidationResult{
			Type:  InputIP,
			Value: ip.String(),
		}, nil
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9._]{3,30}$`)
	if usernameRegex.MatchString(clean) {
		return &ValidationResult{
			Type:  InputUsername,
			Value: clean,
		}, nil
	}

	return nil, errors.New("input is neither a valid IP address nor a supported username")
}



const (
	httpTimeout = 5 * time.Second
	userAgent   = "osint-public-awareness-tool/1.0"
)

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

func (h *HTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	return h.client.Do(req)
}



func ReverseDNS(ip string) []string {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return nil
	}
	return names
}

func IPNetworkType(ip string) string {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "Unknown"
	}
	if parsed.IsPrivate() {
		return "Private IP (Local Network)"
	}
	return "Public IP (ISP / Hosting Provider)"
}


type Platform struct {
	Name string
	URL  string
}

var platforms = []Platform{
	{"GitHub", "https://github.com/%s"},
	{"Twitter/X", "https://x.com/%s"},
	{"Instagram", "https://www.instagram.com/%s"},
	{"Reddit", "https://www.reddit.com/user/%s"},
	{"Medium", "https://medium.com/@%s"},
	{"Dev.to", "https://dev.to/%s"},
}

func CheckUsername(username string) []string {
	client := NewHTTPClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	found := []string{}

	for _, p := range platforms {
		url := fmt.Sprintf(p.URL, username)
		resp, err := client.Get(ctx, url)
		if err != nil {
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			found = append(found, p.Name)
		}
	}

	return found
}



func UsernameExposure(found []string) {
	count := len(found)

	fmt.Println("\n[ Username Exposure ]")
	fmt.Printf("Found on %d platform(s)\n", count)

	if count >= 5 {
		fmt.Println("Exposure Level: HIGH")
	} else if count >= 2 {
		fmt.Println("Exposure Level: MODERATE")
	} else {
		fmt.Println("Exposure Level: LOW")
	}

	if count >= 2 {
		fmt.Println("- Reusing the same username allows strangers to link profiles.")
		fmt.Println("- This can reveal habits, interests, or identity patterns.")
	}
	if count >= 4 {
		fmt.Println("- High reuse increases impersonation and scam risk.")
	}
}

func IPExposure(ip string, rdns []string) {
	fmt.Println("\n[ IP Exposure ]")
	fmt.Println("Network Type:", IPNetworkType(ip))

	if len(rdns) > 0 {
		fmt.Println("Reverse DNS Records Found:")
		for _, h := range rdns {
			fmt.Println(" -", h)
		}
		fmt.Println("Exposure Level: MODERATE")
		fmt.Println("- Reverse DNS can reveal infrastructure patterns.")
	} else {
		fmt.Println("No reverse DNS records found.")
		fmt.Println("Exposure Level: LOW")
	}

	if IPNetworkType(ip) == "Public IP (ISP / Hosting Provider)" {
		fmt.Println("- Public IPs reveal ISP and approximate region.")
		fmt.Println("- Often used in geo-targeted scams and tracking.")
	}
}

func AwarenessFooter() {
	fmt.Println("\n[ Awareness Notice ]")
	fmt.Println("This information is publicly accessible.")
	fmt.Println("No hacking, login, or private access was used.")
	fmt.Println("It means others can see it too.")
}



func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <ip | username>")
		return
	}

	result, err := ValidatePublicInput(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch result.Type {
	case InputIP:
		fmt.Println("Input Type: IP Address")
		rdns := ReverseDNS(result.Value)
		IPExposure(result.Value, rdns)

	case InputUsername:
		fmt.Println("Input Type: Username")
		found := CheckUsername(result.Value)
		UsernameExposure(found)
	}

	AwarenessFooter()
}
