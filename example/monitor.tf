resource "datadog_monitor_metric" "test" {
    name = "App metric"
    message = "This is so good"

    metric = "app.some.metric"
    time_aggr = "avg"
    time_window = "last_5m"
    space_aggr = "avg"
    operator = ">"

    warning {
        threshold = 80
        notify = "@slack-team-infra"
    }

    critical {
        threshold = 100
        notify = "@slack-team-infra @pagerduty"
    }
}