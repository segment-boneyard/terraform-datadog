package monitor

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func Exists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Println("EXISTS")
	client := &http.Client{}
	for _, v := range strings.Split(d.Id(), "__") {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s%s", MONITOR_ENDPOINT, v, AuthSuffix(meta)), nil)
		res, err := client.Do(req)
		if err != nil {
			return false, err
		}
		if res.StatusCode > 400 {
			return false, nil
		}
	}
	return true, nil
}
