
data "changedetection_notifications" "all" {
}

output "example" {
  value = data.changedetection_notifications.all
}
