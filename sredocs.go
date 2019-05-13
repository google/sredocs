package main

import (
	"flag"
	"github.com/google/sredocs/charter"
	"github.com/google/sredocs/exporter/bigquery"
	"github.com/google/sredocs/postmortem"
	"github.com/google/sredocs/source/drive"
	"io/ioutil"
	"log"
)

var (
	downloadCredentials = flag.String("download_credentials_path", "", "Path to credentials for download.")
	downloadFolder      = flag.String("download_folder", "", "Folder to download.")
	downloadDestination = flag.String("download_destination", "", "Path to download to.")
	uploadCredentials   = flag.String("upload_credentials_path", "", "Path to credentials for upload.")
	project             = flag.String("gcp_project", "", "")
	sourceToUpload      = flag.String("source_to_upload", "", "")
	dataset             = flag.String("dataset", "", "")
	table               = flag.String("table", "", "")
	src                 = flag.String("source", "", "Source document to be parsed.")
	kind                = flag.String("kind", "all", "all, charter or postmortem.")
	mode                = flag.String("mode", "parse", "download, parse or upload.")
)

func main() {
	flag.Parse()
	if *kind == "" || *mode == "" {
		log.Fatalf("kind and mode must be set.")
	}
	switch *mode {
	default:
		log.Fatal("unknown mode has been specified.")
	case "download":
		err := drive.Download(*downloadCredentials, *downloadFolder, *downloadDestination)
		if err != nil {
			log.Fatal(err)
		}
	case "upload":
		err := bigquery.Upload(*project, *sourceToUpload, *dataset, *table)
		if err != nil {
			log.Fatal(err)
		}
	case "parse":
		if *src == "" {
			log.Fatalf("source argument must be set in parse mode.")
		}

		b, err := ioutil.ReadFile(*src)
		if err != nil {
			log.Fatal(err)
		}
		switch *kind {
		default:
			fallthrough
		case "charter":
			csv, err := charter.Parse(charter.Fields, b)
			if err != nil {
				log.Println(err)
			}
			log.Println(csv)
			if *kind == "all" {
				goto nextCase
			}
		nextCase:
			fallthrough
		case "postmortem":
			csv, err := postmortem.Parse(postmortem.Fields, b)
			if err != nil {
				log.Println(err)
			}
			log.Println(csv)
		}
	}
}
