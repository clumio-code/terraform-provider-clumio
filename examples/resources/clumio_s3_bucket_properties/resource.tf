resource "clumio_s3_bucket_properties" enable_cb{
  bucket_id = "clumio_assigned_bucket_id"
  event_bridge_enabled = true
}

# Below example shows how to get the bucket_id from the clumio_s3_bucket data source and use
# it in the clumio_s3_bucket_properties resource

data "clumio_s3_bucket" bucket{
  bucket_names = ["test-bucket"]
}

resource "clumio_s3_bucket_properties" enable_cb{
  bucket_id = element(tolist(data.clumio_s3_bucket.bucket.s3_buckets), 0).id
  event_bridge_enabled = true
}
