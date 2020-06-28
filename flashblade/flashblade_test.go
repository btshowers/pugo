package flashblade

import (
	"fmt"
	"net/http"
	"testing"
)

// RoundTripFunc is for returning a test response to the client
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip is the test http Transport
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func testAccGenerateClient(t *testing.T) *Client {

	apiToken := "PUREUSER"
	restVersion := "1.9"
	target := "pureapisim.azurewebsites.net"

	c, err := NewClient(target, apiToken, restVersion)
	if err != nil {
		t.Fatalf("error setting up client: %s", err)
	}
	fmt.Printf(c.XAuthToken + "\n")
	return c
}

func TestAccClient(t *testing.T) {
	//testAccPreChecks(t)
	testAccGenerateClient(t)
}
