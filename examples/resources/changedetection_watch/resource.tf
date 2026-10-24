
resource "changedetection_watch" "localhost" {
  title = "Localhost"
  url   = "https://localhost/"
}

resource "changedetection_tag" "things" {
  title              = "things"
  notification_muted = false
  tag_colour         = "#33ff44"
}

resource "changedetection_watch" "localhost" {
  title               = "Something"
  url                 = "https://example.com/item"
  fetch_backend       = "html_webdriver"
  method              = "POST"
  notification_muted  = true
  processor           = "restock_diff"
  tags                = [changedetection_tag.things.id]
  webdriver_delay     = 3
  paused              = true
  ignore_status_codes = true
  headers = {
    "X-Whatever" = "something"
  }
}
