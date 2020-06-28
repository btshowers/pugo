/*
	Package FlashBlade enables management of Pure FlashBlade using the REST API.package flashblade




*/
package flashblade

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

//Client is the flashblade client
type Client struct {
	Target      string
	APIToken    string
	XAuthToken  string
	RestVersion string
	//UserAgent   string

	client *http.Client

	// Array            *ArrayService
	// Volumes          *VolumeService
	// Hosts            *HostService
	// Hostgroups       *HostgroupService
	// Offloads         *OffloadService
	// Protectiongroups *ProtectiongroupService
	// Vgroups          *VgroupService
	// Networks         *NetworkService
	// Hardware         *HardwareService
	// Users            *UserService
	// Dirsrv           *DirsrvService
	// Pods             *PodService
	// Alerts           *AlertService
	// Messages         *MessageService
	// Snmp             *SnmpService
	// Cert             *CertService
	// SMTP             *SMTPService
}

//NewClient build the client.
func NewClient(target string, apiToken string, restVersion string) (*Client, error) {

	// Get the REST API version to use
	// if restVersion != "" {
	// 	err := checkRestVersion(restVersion, target)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// } else {
	// 	r, err := chooseRestVersion(target)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	restVersion = r
	// }

	// Create a new Client instance
	cookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &Client{Target: target, APIToken: apiToken, RestVersion: restVersion}
	c.client = &http.Client{Transport: tr, Jar: cookieJar}

	err := c.login()
	if err != nil {
		return nil, err
	}

	return c, err
}

//Do execute the client request
func (c *Client) Do(req *http.Request, v interface{}, reestablishSession bool) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println("Do request failed")
		return nil, err
	}
	defer resp.Body.Close()

	// if err := validateResponse(resp); err != nil {
	// 	return resp, err
	// }

	// err = decodeResponse(resp, v)
	return resp, err

}

// func (c *Client) getApiVersion(url) {
// 	//make the rest call
// 	resp := apiCall("GET", "https://"+managementIP.Text()+"/api/api_version", apiToken.Text(), nil)

// 	type Version struct {
// 		Versions []string `json:"versions"`
// 	}

// 	var version Version
// 	err := json.Unmarshal(resp, &version)
// 	if err == nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Printf("%v", (len(version.Versions)))
// 	if len(version.Versions) > 0 {
// 		apiUrlForm.SetText("https://" + managementIP.Text() + "/api/" + version.Versions[(len(version.Versions)-1)])
// 		apiUrlLabel.SetText("https://" + managementIP.Text() + "/api/" + version.Versions[(len(version.Versions)-1)])
// 	}
// 	//set the response in the display of the app
// 	initResult.SetText(string(resp))
// }

//logon and get x-auth-token using api token.
func (c *Client) login() error {
	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	//client := &http.Client{}
	url := ("https://" + c.Target + "/api/login")
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	req.Header.Set("api-token", c.APIToken)
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//set the status code for the response
	// statusCode = resp.StatusCode

	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	//Sets the x-auth-token from the header response
	if len(resp.Header["X-Auth-Token"]) > 0 {
		s := resp.Header["X-Auth-Token"]
		t := strings.Replace(s[0], "[", "", -1)
		t = strings.Replace(t, "]", "", -1)
		c.XAuthToken = t
	}

	return nil
}

// formatPath returns the formated string to be used for the base URL in
// all API calls
func (c *Client) formatPath(path string) string {
	return fmt.Sprintf("https://%s/api/%s/%s", c.Target, c.RestVersion, path)
}
