package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func GetIDFromResponse(h *http.Response) (string, error) {
	body, err := ioutil.ReadAll(h.Body)
	if err != nil {
		return "", err
	}
	h.Body.Close()
	log.Println(h)
	log.Println(string(body))
	v := map[string]interface{}{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		return "", err
	}
	if id, ok := v["id"]; ok {
		return strconv.Itoa(int(id.(float64))), nil
	}
	return "", fmt.Errorf("error getting ID from response %s", h.Status)
}

func MarshalMetric(d *schema.ResourceData, typeStr string) ([]byte, error) {
	name := d.Get("name").(string)
	message := d.Get("message").(string)
	timeAggr := d.Get("time_aggr").(string)
	timeWindow := d.Get("time_window").(string)
	spaceAggr := d.Get("space_aggr").(string)
	metric := d.Get("metric").(string)
	tags := d.Get("metric_tags").(string)
	operator := d.Get("operator").(string)
	var key string
	if k, ok := d.Get("metric_key").(string); ok {
		key = fmt.Sprintf(" by {%s}", k)
	}
	query := fmt.Sprintf("%s(%s):%s:%s{%s}%s %s %s", timeAggr, timeWindow, spaceAggr, metric, tags, key, operator, d.Get(fmt.Sprintf("%s.threshold", typeStr)))

	log.Println(query)
	m := map[string]interface{}{
		"type":    "metric alert",
		"query":   query,
		"name":    fmt.Sprintf("[%s] %s", typeStr, name),
		"message": fmt.Sprintf("%s %s", message, d.Get(fmt.Sprintf("%s.notify", typeStr))),
		"options": map[string]interface{}{
			"notify_no_data":    d.Get("notify_no_data").(bool),
			"no_data_timeframe": d.Get("no_data_timeframe").(int),
		},
	}
	return json.Marshal(m)
}

func AuthSuffix(meta interface{}) string {
	m := meta.(map[string]string)
	return fmt.Sprintf("?api_key=%s&application_key=%s", m["api_key"], m["app_key"])
}

func Create(d *schema.ResourceData, meta interface{}) error {
	warningBody, _ := MarshalMetric(d, "warning")
	criticalBody, _ := MarshalMetric(d, "critical")

	resW, err := http.Post(fmt.Sprintf("%s%s", MONITOR_ENDPOINT, AuthSuffix(meta)), "application/json", bytes.NewReader(warningBody))
	if err != nil {
		return fmt.Errorf("error creating warning: %s", err.Error())
	}

	resC, err := http.Post(fmt.Sprintf("%s%s", MONITOR_ENDPOINT, AuthSuffix(meta)), "application/json", bytes.NewReader(criticalBody))
	if err != nil {
		return fmt.Errorf("error creating critical: %s", err.Error())
	}

	warningMonitorID, err := GetIDFromResponse(resW)
	if err != nil {
		return err
	}
	criticalMonitorID, err := GetIDFromResponse(resC)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s__%s", warningMonitorID, criticalMonitorID))

	return nil
}
