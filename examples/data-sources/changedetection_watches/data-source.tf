data "changedetection_watches" "all" {
}

output "example" {
  value = data.changedetection_watches.all
}
