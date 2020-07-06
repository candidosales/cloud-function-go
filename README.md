# API in GO to create Cloud Function

> An API for deploying Google Cloud Functions in Go.

## 🎯 Features

- Create or update a Cloud Function;
- Get a Cloud Function via function ID;
- Set IAM Policy to enable unauthenticated users invoke the Cloud Function;

## ⚙️ Create a private key to access the GCP

For this method, you'll need to [create a service account](https://cloud.google.com/docs/authentication/getting-started), and download a key.

1. In the GCP Console, go to the [Create service account key](https://console.cloud.google.com/apis/credentials/serviceaccountkey?_ga=2.44822625.-475179053.1491320180) page;
2. From the Service account drop-down list, select New service account;
3. In the Service account name field, enter a name;
4. From the Role drop-down list, select Project > Owner or if you wanna more restrict roles you can choose: `cloudfunctions.admin` (Cloud Functions Admin) and `iam.serviceAccountUser` (Service Account User);
   1. The service account needs to Cloud Functions Admin to have permission to update the IAM policies;
5. Click Create. A JSON file that contains your key downloads to your computer;
6. Save the JSON file to a safe place on your local machine (I create a hidden folder and save it in `~/.google-cloud-keys/key.json` );
7. Update the constant `projectID` with the name of your project at `main.go`;

## Run

```bash
GOOGLE_APPLICATION_CREDENTIALS=~/.google-cloud-keys/key.json go run *.go
```

Send the below request to http://localhost:8080/

## API

<details>
  <summary>Create or update a Cloud Function</summary>

#### Accepts HTTP `POST`

#### Path `/`

#### Request

```json
{
   "name":"cloud-function-test",
   "files":[
      {
         "name":"index.js",
         "content":"exports.helloWorld = (req, res) => {\n  let message = req.query.message || req.body.message || 'Hello World! 4';\n  res.status(200).send(message);\n};"
      },
      {
         "name":"package.json",
         "content":"{\n  \"name\": \"sample-http\",\n  \"version\": \"0.0.1\"\n}"
      }
   ]
}
```

#### Response

```json
{
   "message":"Created the Cloud Function"
}
```

#### Result

![preview](./images/result.png)

![preview](./images/Functions_–_Cloud_Functions_–_Google_Cloud_Platform.png)
</details>
<details>
  <summary>Get a Cloud Function via function ID</summary>

#### Accepts HTTP `GET`

#### Path `/`

#### Request

- `function_id` - string representing the cloud function ID

```bash
curl http://localhost:8080?function_id=trud_test
```

#### Response

```json
{
   "cloudFunction":{
      "availableMemoryMb":128,
      "entryPoint":"helloWorld",
      "httpsTrigger":{
         "url":"https://us-central1-vendasta-hackathon.cloudfunctions.net/trud_test"
      },
      "ingressSettings":"ALLOW_ALL",
      "name":"projects/vendasta-hackathon/locations/us-central1/functions/trud_test",
      "runtime":"nodejs8",
      "serviceAccountEmail":"vendasta-hackathon@appspot.gserviceaccount.com",
      "sourceUploadUrl":"https://storage.googleapis.com/gcf-upload-us-central1-ed9343f1-111c-435e-b16c-28bb09e2f13f/f3679a71-e48e-4c94-8d26-307775296a73.zip?GoogleAccessId=service-481416019804@gcf-admin-robot.iam.gserviceaccount.com&Expires=1590974019&Signature=HDKLbM4x1StRLa8ICqR%2B2R5WEHRmB%2BeqdM2e7btRh1Jb5%2BDqn5SMh1SziOh48ZYDDTjSHpxtQO17%2F3ZeopBXhtEEg2ytGF%2BzJXRA5C2k5BuTX4ULK9OyLQwhE1TDGK5DqJ4JZ%2Fnylfnpdvs4UtqpP4s3Rt4yBszBM5HfIwXYsK7S2HXmskEMc7U9a8rXP6QdqSEwRyTth%2FWTorLt8S9qB0VB8lX8l7xpIsesL0R0mxhfLhn63TQggkx61yOu9qa%2FShjcnqRquOGDEG%2BQ2zEsvKijCaGDZDLq%2Fx0qSvfXSHAQ41uHHnv%2B7Xy6R31j%2Bhe53H3lKLmvUrgL6oM%2Bc5maHQ%3D%3D",
      "status":"ACTIVE",
      "timeout":"60s",
      "updateTime":"2020-06-01T00:44:07.102Z",
      "versionId":"3"
   }
}
```

</details>
<details>
  <summary>Set IAM Policy to `allUsers` have access to invoker the Cloud Function</summary>

#### Accepts HTTP `POST`

#### Path `/policy`

#### Request

```json
{
   "function_id":"trud_test"
}
```

#### Response

```json
{
   "policy":{
      "bindings":[
         {
            "members":[
               "allUsers"
            ],
            "role":"roles/cloudfunctions.invoker"
         }
      ],
      "etag":"BwWm/F4UOXk=",
      "version":1
   }
}
```
</details>

## Libraries

- https://github.com/gofiber/fiber
- https://github.com/awslabs/git-secrets

### Upgrade Go Fiber

```bash
go get -u github.com/gofiber/fiber
```

## 📚 References

- Cloud Function API: https://cloud.google.com/functions/docs/reference/rest
- Client Libraries: https://cloud.google.com/apis/docs/client-libraries-explained
- An API and CLI for deploying Google Cloud Functions in Node.js: https://github.com/JustinBeckwith/gcx
- Deploying from the Google Cloud Functions API: https://cloud.google.com/functions/docs/deploying/api
- Auto-generated Google APIs for Go: https://github.com/googleapis/google-api-go-client
- gcloud functions deploy: https://cloud.google.com/sdk/gcloud/reference/functions/deploy#--runtime
- IAM Policies to Cloud Function: https://cloud.google.com/functions/docs/reference/iam/roles
- Help keep your Google Cloud service account keys safe: taking charge of your security: https://cloud.google.com/blog/products/gcp/help-keep-your-google-cloud-service-account-keys-safe

## Author

- Cândido Sales - [@candidosales](https://twitter.com/candidosales)
