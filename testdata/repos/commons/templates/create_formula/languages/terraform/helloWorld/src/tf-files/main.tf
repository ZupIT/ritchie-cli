provider "aws" {
  region = var.bucket_region
}

resource "aws_s3_bucket" "storage-s3" {
  bucket = var.bucket_name
}
