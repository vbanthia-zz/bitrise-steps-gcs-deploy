package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"

	"cloud.google.com/go/storage"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func downloadFile(downloadURL, targetPath string) error {
	outFile, err := os.Create(targetPath)
	if err != nil {
		failf("Failed to create (%s), error: %s", targetPath, err)
	}
	defer func() {
		if err = outFile.Close(); err != nil {
			log.Warnf("Failed to close (%s), error: %s", targetPath, err)
		}
	}()

	resp, err := http.Get(downloadURL)
	if err != nil {
		failf("Failed to download from (%s), error: %s", downloadURL, err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Warnf("Failed to close (%s) body", downloadURL)
		}
	}()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		failf("Failed to download from (%s), error: %s", downloadURL, err)
	}

	return nil
}

func main() {
	keyPath := os.Getenv("service_account_json_key_path")
	projectID := os.Getenv("project_id")
	bucketName := os.Getenv("bucket_name")
	folderName := os.Getenv("folder_name")
	uploadFilePath := os.Getenv("upload_file_path")
	uploadedFileName := os.Getenv("uploaded_file_name")

	// Download json_key if bitrise file storage is used
	if strings.HasPrefix(keyPath, "http") {
		tmpDir, err := pathutil.NormalizedOSTempDirPath("__google-cloud-storage__")

		if err != nil {
			failf("Failed to create tmp dir, error: %s", err)
		}

		targetPath := filepath.Join(tmpDir, "key.json")
		if err := downloadFile(keyPath, targetPath); err != nil {
			failf("Failed to download json key file, error: %s", err)
		}

		keyPath = targetPath
	}

	// Set GOOGLE_APPLICATION_CREDENTIALS to json_key so that
	// gcloud library uses service account
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", keyPath)

	// Creat client
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithServiceAccountFile(keyPath))

	if err != nil {
		failf("Failed to create new storage client, error: %s", err)
	}

	// Create bucket if it does not exist
	bucketExist := false
	it := client.Buckets(ctx, projectID)

	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			failf("Failed to create new storage client, error: %s", err)
		}

		if battrs.Name == bucketName {
			bucketExist = true
			log.Infof("Bucket %s already exist", bucketName)
		}
	}

	if !bucketExist {
		if err = client.Bucket(bucketName).Create(ctx, projectID, nil); err != nil {
			failf("Failed to create bucket, error: %s", err)
		}
		log.Infof("Bucket %s created successfully", bucketName)
	}

	// Uploading file
	file, err := os.Open(uploadFilePath)
	if err != nil {
		failf("File (%s) does not exist, error: %s", uploadFilePath, err)
	}

	defer file.Close()

	if folderName != "" {
		uploadedFileName = folderName + "/" + uploadedFileName
	}

	wc := client.Bucket(bucketName).Object(uploadedFileName).NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		failf("File (%s) does not exist, error: %s", uploadFilePath, err)
	}

	if err = wc.Close(); err != nil {
		failf("Failed to close wc, error: %s", err)
	}

	if err = client.Close(); err != nil {
		failf("Failed to close storage client, error: %s", err)
	}
}
