resource "hcs_obs_bucket" "mywebsite" {
  bucket = "mywebsite"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}

# granting the Read-Only permission to anonymous users
resource "hcs_obs_bucket_policy" "policy" {
  bucket = hcs_obs_bucket.mywebsite.bucket
  policy = <<POLICY
{
  "Statement": [
    {
      "Sid": "AddPerm",
      "Effect": "Allow",
      "Principal": {"ID": "*"},
      "Action": ["GetObject"],
      "Resource": "mywebsite/*"
    } 
  ]
}
POLICY
}

# put index.html
resource "hcs_obs_bucket_object" "index" {
  bucket = hcs_obs_bucket.mywebsite.bucket
  key    = "index.html"
  source = "index.html"
}

# put error.html
resource "hcs_obs_bucket_object" "error" {
  bucket = hcs_obs_bucket.mywebsite.bucket
  key    = "error.html"
  source = "error.html"
}
