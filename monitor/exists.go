package monitor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func Exists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	b = true
	for _, v := range strings.Split(d.Id(), "__") {
		res, err := http.Get(fmt.Sprintf("%s/%s%s", MONITOR_ENDPOINT, v, AuthSuffix(meta)))
		if err != nil {
			e = err
			continue
		}
		if res.StatusCode > 400 {
			b = false
			continue
		}
		b = b && true
	}
	if !b {
		Delete(d, meta)
	}
	return
}
