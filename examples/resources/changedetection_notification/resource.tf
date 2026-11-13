
resource "changedetection_notification" "all" {
  notification_urls = [
    "get://127.0.0.1/e",
    "get://127.0.0.1/a",
    "get://127.0.0.1/b",
    "get://127.0.0.1/c",
    "get://127.0.0.1/d",
  ]
}

