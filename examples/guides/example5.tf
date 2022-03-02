provider "aws" {
  region     = "us-west-2"
  access_key = "my-access-key"
  secret_key = "my-secret-key"

  # If a session token is required ...
  token = "my-token"
}
