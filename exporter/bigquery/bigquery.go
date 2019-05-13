package bigquery

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
)

func Upload(project, sourcePath, dataset, table string) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, project)
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
		f, err := os.Open(sp.Name())
		if err != nil {
			return err
		}
		source := bigquery.NewReaderSource(f)
		source.AutoDetect = true   // Allow BigQuery to determine schema.
		source.SkipLeadingRows = 1 // CSV has a single header line.
		// TODO(stratus): Add WRITE_TRUNCATE

		loader := client.Dataset(dataset).Table(tableName).LoaderFrom(source)

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
