terraform {
  backend "s3" {
    bucket = "gas-express-tf-prd"
    key    = "terraform.tfstate"
    region = "us-east-1"
  }
}
