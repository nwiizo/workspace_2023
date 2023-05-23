locals {
  list = [
    "foo",
    "bar",
    "baz",
  ]
}

output "output_list" {
  value = [for s in local.list : upper(s) if s != "bar"]
}
