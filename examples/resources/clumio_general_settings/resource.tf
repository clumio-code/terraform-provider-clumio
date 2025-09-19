data "http" "current_ip" {
  url = "https://ipv4.icanhazip.com"
}

resource "clumio_general_settings" "example" {
  auto_logout_duration         = 1200    // 20 minutes
  password_expiration_duration = 7776000 // 90 days
  ip_allowlist                 = ["${chomp(data.http.current_ip.response_body)}/32",]
}
