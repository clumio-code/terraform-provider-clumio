data clumio_protection_group_asset example {
  protection_group_id = data.clumio_protection_group.example.id
  bucket_id = tolist(data.clumio_s3_bucket.example.s3_buckets)[0].id
}

data "clumio_protection_group" "example" {
  name = "protection-group-name"
}

data "clumio_s3_bucket" "example" {
  bucket_names=["bucket1"]
}
