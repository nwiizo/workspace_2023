locals {
  list = [
    "foo",
    "bar",
    "baz",
  ]
}



output "output_list" {
  value = [for s in local.list : upper(s)]
}

output "output_count" {
  value = length(local.list)
}

output "name_count" {
  value = local.list[length(local.list) - 1]
}
