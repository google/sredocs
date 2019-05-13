package drive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func Download(credentials_path string, folder string, destination string) error {
	b, err := ioutil.ReadFile(credentials_path)
	if err != nil {
		return err
	}
	// TODO(stratus): Eval if a downgrade to drive.DriveFileScope is feasible.
	config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.New(client)
	if err != nil {
		return err
	}
	r, err := srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		return err
	}
	var errorCount int
	if len(r.Files) == 0 {
		return fmt.Errorf("No files found in %s.", folder)
	} else {
		for _, i := range r.Files {
			log.Printf("drive: %s (%s)\n", i.Name, i.Id)
			res, err := srv.Files.Export(i.Id, "text/plain").Download()
			if err != nil {
				errorCount++
				log.Printf("%q while downloading %s (%s)", err, i.Name, i.Id)
			}
			destFile := fmt.Sprintf("%s/%s", destination, i.Name)
			s, err := ioutil.ReadAll(res.Body)
			if err != nil {
				errorCount++
				log.Printf("%q while reading %s", err)
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
