package bigquery

import (
	"cloud.google.com/go/bigquery"
	"context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io/ioutil"
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
	return false
}

func Upload(credentials_path, project, sourcePath, datasetName, table string) error {
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
	// TODO(stratus): Turn per file errors non-fatal.
	for _, sp := range files {
		f, err := os.Open(filepath.Join(sourcePath, sp.Name()))
		if err != nil {
			return err
		}
		source := bigquery.NewReaderSource(f)
		source.AutoDetect = true   // Allow BigQuery to determine schema.
		source.SkipLeadingRows = 1 // CSV has a single header line.
		// TODO(stratus): Add WRITE_TRUNCATE

		d := &dataset{ctx: ctx, client: client, name: datasetName, location: "US"}
		if !d.exists() {
			err := d.create()
			if err != nil {
				return err
			}
		}
		loader := client.Dataset(datasetName).Table(tableName).LoaderFrom(source)

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
