# supermarket-api

The supermarket-api is a Golang application that responds to HTTP requests for updating produce information in a backend database. The application is meant to be deployed in docker.

## Deploy Instructions

To quickly run locally, you can use:

```bash
docker run -d -p 6620:6620 dirtysoc/supermarket-api
```

## Data Types

**Produce**:

```json
{
  "name":  "produce name",
  "produceCode": "16 character alphanumeric identifier with dashes every 4 characters",
  "unitPrice": "The unit price of the produce in USD with 2 digit accuracy",
}
```

## Endpoints

- [x] `GET /produce` returns all produce as JSON
- [x] `GET /produce/{id}` returns details of a specific produce as JSON
- [x] `POST /produce` add a produce item to the database*
- [x] `POST /produce` add multiple produce items to the database
- [x] `DELETE /produce/{id}` removes produce from the database by produceCode

\* Note that all JSON body data must be a JSON array. For example, adding a single new produce item requires that the JSON body of the POST request be a JSON array with a single object in it.

## CI/CD

Builds are automated and all pushes to the master branch are tested, built and uploaded to Docker Hub.

## Local Development

1. Install Golang
2. Clone repo
3. `cd` into repo
4. Edit and run `go test` or build with `go build` and test.

## User Stories

| User Stories                | Narrative                                                                                 |
|-----------------------------|-------------------------------------------------------------------------------------------|
| Adding new produce          | As an employee, I want to add produce, so that I can add items to the database            |
| Deleting a produce item     | As an employee, I want to delete produce, so I can remove produce from the database       |
| Fetch the produce inventory | As an employee, I want to look up produce, so that I understand what produce is available |
