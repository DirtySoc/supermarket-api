# supermarket-api

The supermarket-api is a Golang application that responds to HTTP requests for updating produce information in a backend database. The application is meant to deployed in docker.

## Local Development

1. Install Golang
2. Clone repo
3. `cd` into repo
4. Edit and run `go test` or build with `go build` and test.

## Data Types

```json
{
  "name":  "produce name",
  "produceCode": "16 character alphanumeric identifier with dashes every 4 characters",
  "unitPrice": "The unit price of the produce in USD with 2 digit accuracy",
}
```

## Endpoints

- [ ] `GET /produce` returns all produce as JSON
- [ ] `GET /produce/{id}` returns details of a specific produce as JSON
- [ ] `POST /produce/{id}` add a produce item to the database
- [ ] `DELETE /produce/{id} removed produce from the database by produceCode

## User Stories

| User Stories                | Narrative                                                                                 |
|-----------------------------|-------------------------------------------------------------------------------------------|
| Adding new produce          | As an employee, I want to add produce, so that I can add items to the database            |
| Deleting a produce item     | As an employee, I want to delete produce, so I can remove produce from the database       |
| Fetch the produce inventory | As an employee, I want to look up produce, so that I understand what produce is available |
