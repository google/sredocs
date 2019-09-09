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

package drive

import (
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func search(srv *drive.Service, query string) ([]*drive.File, error) {
	r, err := srv.Files.List().Spaces("drive").Corpora("user").Q(query).Do()
	if err != nil {
		return nil, err
	}
	return r.Files, err
}

func Download(credentials_path string, folder string, destination string) error {
	b, err := ioutil.ReadFile(credentials_path)
	if err != nil {
		return err
	}
	config, err := google.JWTConfigFromJSON(b, drive.DriveReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v.", err)
	}
	srv, err := drive.New(config.Client(oauth2.NoContext))
	if err != nil {
		return err
	}

	q := fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.folder'", folder)
	fmt.Println(q)
	d, err := search(srv, q)
	if err != nil {
		return err
	}
	if len(d) != 1 {
		return fmt.Errorf("Expected 1 match for %s, got %d.", folder, len(d))
	}

	q = fmt.Sprintf("'%s' in parents", d[0].Id)
	files, err := search(srv, q)
	if err != nil {
		return err
	}

	var errorCount int
	if len(files) == 0 {
		return fmt.Errorf("No files found in %s.", folder)
	} else {
		for _, f := range files {
			log.Printf("drive: %s (%s)\n", f.Name, f.Id)
			res, err := srv.Files.Export(f.Id, "text/plain").Download()
			if err != nil {
				errorCount++
				log.Printf("%q while downloading %s (%s)", err, f.Name, f.Id)
			}
			destFile := fmt.Sprintf("%s/%s", destination, f.Name)
			s, err := ioutil.ReadAll(res.Body)
			if err != nil {
				errorCount++
				log.Printf("%q while reading %s", err, f.Name)
			}
			err = ioutil.WriteFile(destFile, s, 0644)
			if err != nil {
				errorCount++
				log.Printf("%q while writing %s", err, destFile)
			}
		}
	}
	if errorCount != 0 {
		return fmt.Errorf("%d errors while downloading files", errorCount)
	}
	return nil
}
