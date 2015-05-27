package monitor

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func Update(d *schema.ResourceData, meta interface{}) error {
	split := strings.Split(d.Id(), "__")
	warningID, criticalID := split[0], split[1]

	warningBody, _ := MarshalMetric(d, "warning")
	criticalBody, _ := MarshalMetric(d, "critical")

	client := http.Client{}

	reqW, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s%s", MONITOR_ENDPOINT, warningID, AuthSuffix(meta)), bytes.NewReader(warningBody))
	resW, err := client.Do(reqW)
	if err != nil {
		return fmt.Errorf("error updating warning: %s", err.Error())
	}
	resW.Body.Close()
	if resW.StatusCode > 400 {
		return fmt.Errorf("error updating warning monitor: %s", resW.Status)
	}

	reqC, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s%s", MONITOR_ENDPOINT, criticalID, AuthSuffix(meta)), bytes.NewReader(criticalBody))
	resC, err := client.Do(reqC)
	if err != nil {
		return fmt.Errorf("error updating critical: %s", err.Error())
	}
	resW.Body.Close()
	if resW.StatusCode > 400 {
		return fmt.Errorf("error updating critical monitor: %s", resC.Status)
	}
	return nil
}
