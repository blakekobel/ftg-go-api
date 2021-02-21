# Blake Kobel FTG Go Skills Challenge
API Gateway link is: {Insert Link}

## Setup 
- Install Terraform here: {Insert Terraform Link}
- Put Terraform.exe file in folder within path
- Create GCP Project
- Enable the following Services/APIs: 
  - Cloud Function, Storage, Cloud Build, API Gateway
- Go to IAM and create a service account
- Give that Service Account Owner owner permissions for the project
- Create a JSON key-pair and save the file that is downloaded

## Resource Creation
- Put JSON key-pair file into this directory (gitignore checks for JSON files)
- Navigate to project directory using CMD or Git Bash
- Edit main.tf lines 10-25
  - put relative file path in the credential file area
  - put project id in the project lines
- Edit openapi.yml line 14
  - https://us-central1-{project-id}.cloudfunctions.net/ftg-ip-api-get
  - Line 14 is what connects the API configuration to the cloud function
- Run the command: Terraform init
  - This command initializes the terraform project and will create a .terraform folder
  - Fix any errors that the initialization may request you do
- Run the command: Terraform apply
  - This command builds all the needed infrastructure for your process

## Review and Consumption
- Log into GCP and see all the resources created in your project
- Go to API Gateway within GCP to determine your API URL
- Begin using your API URL

## Cleanup
- Get back into the Project Directory
- Run the command: Terraform destroy
  - This deletes the project within the directory