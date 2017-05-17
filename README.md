## Bitrise GCS Upload Step

Bitrise Step for uploading build artifacts to Google Cloud Storage.

### How to use it.

#### Upload Service Account JSON File
Bitrise GCS Upload Step uses service account json file for authentication. If you don't have service account yet then first create one using [these instructions](https://cloud.google.com/docs/authentication).

Once you have credentials in JSON file, upload them to Bitrise in `Workflow -> Code Signing -> Generic File Storage`

Bitrise GCS Upload Step uses `BITRISEIO_GCS_SERVICE_ACCOUNT_JSON_FILE_URL` as default key. In case if you choose different key name, you will have to write it in step inputs.

#### Update `bitrise.yml`
Add following step into your `bitrise.yml`

```yaml
- git::https://github.com/vbanthia/bitrise-steps-gcs-deploy.git@1.0:
    title: Upload Apk to Google Cloud Storage
    inputs:
    - service_account_json_key_path: "$BITRISEIO_GCS_SERVICE_ACCOUNT_JSON_FILE_URL"
    - project_id: 'project-name'
    - bucket_name: 'bucket-name'
    - folder_name: 'folder1/folder2'
    - upload_file_path: "$BITRISE_APK_PATH"
    - uploaded_file_name: my-app-$BITRISE_BUILD_NUMBER.apk
```
