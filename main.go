package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <url>")
		return
	}

	targetURL := os.Args[1]

	// Parse the target URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	redirectParams := []string{"redirectUrl"}

	// Create a redirect URL with the target host and path
	redirectURL := "https://example.com"

	// Create a query string with a redirect URL parameter
	for _, redirectParam := range redirectParams {
		query := parsedURL.Query()
		query.Set(redirectParam, redirectURL)
		parsedURL.RawQuery = query.Encode()

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		// Send a GET request to the target URL with the redirect parameter
		resp, err := client.Get(parsedURL.String())
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		// Check if the response redirected to a different host
		if resp.StatusCode >= 300 && resp.StatusCode <= 399 {
			location, err := resp.Location()
			if err != nil {
				fmt.Println("Redirect location not found:", err)
				return
			}

			if location.String() == redirectURL {
				fmt.Println("Vulnerable to open redirect!")
				break
			}
		}
	}
}
