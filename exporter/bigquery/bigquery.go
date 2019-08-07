// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bigquery

import (
	"cloud.google.com/go/bigquery"
	"context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type dataset struct {
	ctx      context.Context
	client   *bigquery.Client
	name     string
	location string
}

func (d *dataset) create() error {
	meta := &bigquery.DatasetMetadata{
		Location: d.location,
	}
	if err := d.client.Dataset(d.name).Create(d.ctx, meta); err != nil {
		return err
	}
	return nil
}

func (d *dataset) exists() bool {
	it := d.client.Datasets(d.ctx)
	for {
		dataset, err := it.Next()
		if err == iterator.Done {
			return false
		}
		if dataset.DatasetID == d.name {
			return true
		}
	}
}

func Upload(credentials_path, project, sourcePath, datasetName, table string, truncate bool) error {
	ctx := context.Background()
	opts := option.WithCredentialsFile(credentials_path)
	client, err := bigquery.NewClient(ctx, project, opts)
	if err != nil {
		return err
	}

	tableName := table
	if tableName == "" {
		tableName = time.Now().Format("20060102")
	}

	files, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		return err
	}
	for i, sp := range files {
		fn := filepath.Join(sourcePath, sp.Name())
		f, err := os.Open(fn)
		if err != nil {
			if len(files) > 1 {
				log.Printf("Error %q while processing %s. Continuing...", err, fn)
			} else {
				return err
			}
		}
		source := bigquery.NewReaderSource(f)
		source.AutoDetect = true   // Allow BigQuery to determine schema.
		source.SkipLeadingRows = 1 // CSV has a single header line.

		d := &dataset{ctx: ctx, client: client, name: datasetName, location: "US"}
		if !d.exists() {
			log.Printf("Creating %s in %s", datasetName, project)
			err := d.create()
			if err != nil {
				return err
			}
		}
		loader := client.Dataset(datasetName).Table(tableName).LoaderFrom(source)
		if truncate && i == 0 {
			log.Printf("Truncate has been set, will override BigQuery table existing data (if any).")
			loader.WriteDisposition = bigquery.WriteTruncate
		}
		log.Printf("Uploading %s to %s (%s).", fn, datasetName, tableName)
		job, err := loader.Run(ctx)
		if err != nil {
			return err
		}
		status, err := job.Wait(ctx)
		if err != nil {
			return err
		}
		if err := status.Err(); err != nil {
			return err
		}
	}
	return nil
}
