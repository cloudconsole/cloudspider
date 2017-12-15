package crawler

import (
	"fmt"
	"io/ioutil"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

func main() {
	config, _ := edgegrid.Init("~/.edgerc", "default")

	// Retrieve dig information for specified location
	req, _ := client.NewRequest(config, "GET", "/diagnostic-tools/v1/dig", nil)

	q := req.URL.Query()
	q.Add("hostname", "developer.akamai.com")
	q.Add("queryType", "A")
	q.Add("location", "Auckland, New Zealand")

	req.URL.RawQuery = q.Encode()
	resp, _ := client.Do(config, req)

	defer resp.Body.Close()
	byt, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byt))
}
