# sredocs

This repository contains templates for SRE related documents, a parser
and related tools.

The documents are both human and machine readable.


## Download

Here is how to make the download mode work properly:

Go to https://console.developers.google.com.
I'd recommend creating a new project called "sredocs".
Switch to that project.
Search for "Google Drive API" then enable it.
Go to APIs & Services -> Credentials.
Click on "Create credentials" then select Service account key.
Give your Service account key a name, e.g "sredocs".
Take note of its address (looks like an email address).
Do not assign it a role while creating it.
Save the private key (it will be automatically downloaded).
Switch to Google Drive.
Give your service account address, e.g sredocs@...iam.gservice... view access to the folder with SRE docs.
Pass the private key to sredocs via -download_credentials_path


## Upload

Follow the instructions for Download first.
Go to https://console.developers.google.com.
Navigate to IAM & Admin
Click on Add
Search for your service account, e.g sredocs@...
Give it the BigQuery Admin role
