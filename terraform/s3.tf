resource "aws_s3_bucket" "gas_express_s3_bucket" { // creates a new `aws_s3_bucket` resource and gives it an id `gas_express_s3_bucket`
  bucket        = var.bucket                       // names the bucket `gas-express`
  acl           = "public-read"                    // defines the bucket as private
  force_destroy = true

  website {
    index_document = "gasExpress.json"
    error_document = "tmp/gasExpress.json"
  }

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET"]
    allowed_origins = ["*"]
    max_age_seconds = 1800
  }

  versioning {
    enabled = true // enables versioning for the `gas_express_s3_bucket` resource
  }
  policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "PublicReadForGetBucketObjects",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::${var.bucket}/*"
    }
  ]
}
EOF

  tags = {
    Name = "gas-express" // sets bucket tag in aws to "gas-express"
  }
}
