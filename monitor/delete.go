package monitor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func Delete(d *schema.ResourceData, meta interface{}) (e error) {
	for _, v := range strings.Split(d.Id(), "__") {
		client := http.Client{}
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s%s", MONITOR_ENDPOINT, v, AuthSuffix(meta)), nil)
		_, err := client.Do(req)
		e = err
	}
	return
}
