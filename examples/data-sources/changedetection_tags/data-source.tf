data "changedetection_tags" "all" {
}

output "example" {
  value = data.changedetection_tags.all
}
