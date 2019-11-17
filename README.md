# sredocs

sredocs is to turn postmortems, charters and other SRE docs into machine readable data for analysis while retaining human readability properties.

This repository contains code for downloading postmortems and/or charters following certain template(s) from Google Docs, parsing the data to CSV and support for uploading the result to Google BigQuery.

Here is a [demo dashboard for incident analysis using this, Google Data Studio and BigQuery](http://datastudio.google.com/reporting/81474668-0f5e-409b-9508-51bd38a5321c). It contains an interactive report with information related to five dummy postmortems.

sredocs currently parses postmortems using a single specific template [(sample)](https://docs.google.com/document/d/15Vyd2b3RYQZ3p6UztfbHWsd0fCzzwa42jwEczwFB8es/preview) and charter that will be published separately later.

## Download

Here is how to use the download mode:

1. Go to https://console.developers.google.com. I'd recommend creating a new project called "sredocs". 
1. Switch to that project. Search for "Google Drive API" then enable it. 
1. Go to APIs & Services -> Credentials. 
1. Click on "Create credentials" then select Service account key. 
1. Give your Service account key a name, e.g "sredocs". 
1. Take note of its address (looks like an email address). Do not assign it a role while creating it. 
1. Save the private key (it will be automatically downloaded). 
1. Switch to Google Drive. 
1. Give your service account address, e.g sredocs@...iam.gservice... view access to the folder with SRE docs.
1. Pass the private key to sredocs via -download_credentials_path.
1. Run `mkdir <download_dir>`
1. Run `sredocs -mode=download -cloud_credentials=<creds.json> -download_folder=<docs_folder> -download_destination=<download_dir>`

## Parse

Here is how to use the parse mode:

1. Run `mkdir <parsed_dir>`
1. Run `sredocs -mode=parse -parse_kind=auto -parse_path=<download_dir> -parse_output_path=<parsed_dir>`

## Upload

Here is how to use the upload mode:

1. Follow the instructions for Download first. 
1. Go to https://console.developers.google.com. 
1. Navigate to IAM & Admin. 
1. Click on the Add button. 
1. Search for your service account, e.g sredocs@... Give it the BigQuery Admin role.
1. Run `sredocs -mode=upload -cloud_credentials=<creds.json> -upload_path=<parsed_dir> -upload_project=<org:project> -upload_dataset=<sredocs> -upload_table=<sredocs_20190603> -upload_truncate=true`

## BYOP

You can also bring your own parser.

~~~~
type Parser interface {
        CompileRegex(fields []string) ([]*regexp.Regexp, error)
        Parse(fields []string, b []byte) (*bytes.Buffer, error)
        CSVHeader(regexps []*regexp.Regexp) []string
        NamedGroup(field string) string
        Save(buf *bytes.Buffer, filename string) error
}
~~~~

You will have to modify charter/* and postmortem/* to switch away from the DefaultParser.

You can also extend sredocs by adding a new kind of document by using charter/* as a template.
