terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "3.5.0"
    }
  }
}

provider "google" {

  credentials = file("axiomatic-set-305500-5610f6d7ec1a.json")

  project = "axiomatic-set-305500"
  region  = "us-central1"
  zone    = "us-central1-c"
}

provider "google-beta" {

  credentials = file("axiomatic-set-305500-5610f6d7ec1a.json")

  project = "axiomatic-set-305500"
  region  = "us-central1"
}

resource "google_storage_bucket" "ftg_ip_api_bucket" {
  name = "ftg-ipapi-bucket"
}

data "archive_file" "ip_api" {
  type          = "zip"
  output_path   = "${path.module}/files/ip_api.zip"
  source_dir    = "${path.root}/ip_api/"
}

resource "google_storage_bucket_object" "archive" {
  name          = "ip_api.zip"
  bucket        = google_storage_bucket.ftg_ip_api_bucket.name
  source        = "${path.module}/files/ip_api.zip"
  depends_on    = [data.archive_file.ip_api]
}

#Create the GoLang Cloud Function
resource "google_cloudfunctions_function" "test" {
    name                      = "ftg-api-ip-get"
    entry_point               = "PrintIP"
    available_memory_mb       = 256
    timeout                   = 120
    region                    = "us-central1"
    trigger_http              = true
    # ingress_settings              = "ALLOW_ALL"
    runtime                   = "go113"
    source_archive_bucket     = google_storage_bucket.ftg_ip_api_bucket.name
    source_archive_object     = google_storage_bucket_object.archive.name
    labels                    =  {"deployment_name":"test"}
}

# IAM entry for all users to invoke the function
resource "google_cloudfunctions_function_iam_member" "invoker" {
  project        = google_cloudfunctions_function.test.project
  region         = google_cloudfunctions_function.test.region
  cloud_function = google_cloudfunctions_function.test.name

  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}

resource "google_api_gateway_api" "FTG_api_gateway_api" {
  provider = google-beta
  api_id = "ftg-api-ip"
}

resource "google_api_gateway_api_config" "FTG_api_gateway_cfg" {
  provider = google-beta
  api = google_api_gateway_api.FTG_api_gateway_api.api_id
  openapi_documents {
    document {
        path = "openapi.yml"
        contents = filebase64("openapi.yml")
    }
  }
}

resource "google_api_gateway_gateway" "FTG_api_gateway" {
  provider = google-beta
  api_config = google_api_gateway_api_config.FTG_api_gateway_cfg.id
  gateway_id = "ftg-api-ip-gateway"
}