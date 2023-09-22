check "health_check" {
  data "http" "example" {
    url = "https://blog.3-shake.com/"
  }
 
  assert {
    condition     = data.http.example.status_code == 200
    error_message = "blog.3-shake.com returned an unhealthy status code"
  }
}
