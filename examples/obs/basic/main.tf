resource "hcs_obs_bucket" "myexample" {
  bucket = "myexample-bucket"
  acl    = "private"

  tags = {
    type = "bucket"
    env  = "Test"
  }
}

# put myobject1 by content
resource "hcs_obs_bucket_object" "myobject1" {
  bucket       = hcs_obs_bucket.myexample.bucket
  key          = "myobject1"
  content      = "content of myobject1"
  content_type = "application/xml"
}

# put myobject2 by source file
resource "hcs_obs_bucket_object" "myobject2" {
  bucket = hcs_obs_bucket.myexample.bucket
  key    = "myobject2"
  source = "hello.txt"
}

# put myobject3 by source file and encryption with default key
resource "hcs_obs_bucket_object" "myobject3" {
  bucket     = hcs_obs_bucket.myexample.bucket
  key        = "myobject3"
  source     = "hello.txt"
  encryption = true
}
