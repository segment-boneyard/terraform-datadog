package monitor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func Exists(d *schema.ResourceData, meta interface{}) (bool, error) {
	for _, v := range strings.Split(d.Id(), "__") {
		res, err := http.Get(fmt.Sprintf("%s/%s%s", MONITOR_ENDPOINT, v, AuthSuffix(meta)))
		if err != nil {
			return false, err
		}
		if res.StatusCode > 400 {
			return false, nil
		}
	}
	return true, nil
}
