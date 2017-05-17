package main

import (
	"fmt"
	"log"
	"io"
	"os"
	"net/http"
	"path/filepath"
	"strings"

	"google.golang.org/api/storage/v1"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)

func downloadFile(downloadURL, targetPath string) error {
	outFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create (%s), error: %s", targetPath, err)
	}
	defer func() {
		if err := outFile.Close(); err != nil {
			log.Warnf("Failed to close (%s)", targetPath)
		}
	}()

	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download from (%s), error: %s", downloadURL, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Warnf("failed to close (%s) body", downloadURL)
		}
	}()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to download from (%s), error: %s", downloadURL, err)
	}

	return nil
}


// ListBuckets returns a slice of all the buckets for a given project.
func ListBuckets(projectID string) ([]*storage.Bucket, error) {
	ctx := context.Background()

	// Create the client that uses Application Default Credentials
	// See https://developers.google.com/identity/protocols/application-default-credentials
	client, err := google.DefaultClient(ctx, storage.DevstorageReadOnlyScope)
	if err != nil {
		return nil, err
	}

	// Create the Google Cloud Storage service
	service, err := storage.New(client)
	if err != nil {
		return nil, err
	}

	buckets, err := service.Buckets.List(projectID).Do()
	if err != nil {
		return nil, err
	}

	return buckets.Items, nil
}

func main() {
	key_path := os.Getenv("service_account_json_key_path")

	if strings.HasPrefix(key_path, "http") {
		downloadUrl := key_path
		tmpDir, err := pathutil.NormalizedOSTempDirPath("__google-cloud-storage__")

		if err != nil {
		  fmt.Errorf("Failed to create tmp dir, error: %s", err)
		}

		targetPath := filepath.Join(tmpDir, "key.json")
		if err := downloadFile(key_path, targetPath); err != nil {
				fmt.Errorf("Failed to download json key file, error: %s", err)
		}

		key_path = targetPath
	}

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", key_path)

	if len(os.Args) < 2 {
		fmt.Println("usage: listbuckets <projectID>")
		os.Exit(1)
	}
	project := os.Args[1]

	buckets, err := ListBuckets(project)
	if err != nil {
		log.Fatal(err)
	}

	// Print out the results
	for _, bucket := range buckets {
		fmt.Println(bucket.Name)
	}
}
