# Data source to fetch the bucket details based on bucket name(s).
data "clumio_s3_bucket" "example" {
  bucket_names=["bucket1", "bucket2"]
}

# Use the data source to determine the bucket id.
resource "clumio_protection_group_bucket" "advanced_example"{
  count = length(data.clumio_s3_bucket.example.s3_buckets)
  protection_group_id = "protection-group-id"
  bucket_id = tolist(data.clumio_s3_bucket.example.s3_buckets)[count.index].id
}
