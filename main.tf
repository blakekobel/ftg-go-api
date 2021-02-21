terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "3.5.0"
    }
  }
}

provider "google" {

  credentials = file("formal-theater-305519-c85df13a9aba.json")

  project = "formal-theater-305519"
  region  = "us-central1"
  zone    = "us-central1-c"
}

provider "google-beta" {

  credentials = file("formal-theater-305519-c85df13a9aba.json")

  project = "formal-theater-305519"
  region  = "us-central1"
}

resource "google_storage_bucket" "ftgapiipbucket" {
  name = "ftgapiipbucket"
}

data "archive_file" "api_ip" {
  type          = "zip"
  output_path   = "${path.module}/files/api_ip.zip"
  source_dir    = "${path.root}/ip_api/"
}

resource "google_storage_bucket_object" "archive" {
  name          = "api_ip.zip"
  bucket        = google_storage_bucket.ftgapiipbucket.name
  source        = "${path.module}/files/api_ip.zip"
  depends_on    = [data.archive_file.api_ip]
}

#Create the GoLang Cloud Function
resource "google_cloudfunctions_function" "test" {
    name                      = "ftg-ip-api-get"
    entry_point               = "PrintIP"
    available_memory_mb       = 256
    timeout                   = 120
    region                    = "us-central1"
    trigger_http              = true
    # ingress_settings              = "ALLOW_ALL"
    runtime                   = "go113"
    source_archive_bucket     = google_storage_bucket.ftgapiipbucket.name
    source_archive_object     = google_storage_bucket_object.archive.name
    labels                    =  {"deployment_name":"test"}
}

##Uncomment if the API Gateway cant call the function or if you want to test the function by itself.
# # IAM entry for all users to invoke the function
# resource "google_cloudfunctions_function_iam_member" "invoker" {
#   project        = google_cloudfunctions_function.test.project
#   region         = google_cloudfunctions_function.test.region
#   cloud_function = google_cloudfunctions_function.test.name

#   role   = "roles/cloudfunctions.invoker"
#   member = "allUsers"
# }

resource "google_api_gateway_api" "FTG_api_gateway_api_ip" {
  provider = google-beta
  api_id = "ftg-api-gw"
}

resource "google_api_gateway_api_config" "FTG_api_gateway_cfg" {
  provider = google-beta
  api = google_api_gateway_api.FTG_api_gateway_api_ip.api_id
  display_name = "ftg-api-ip-cfg"
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