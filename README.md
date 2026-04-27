# LoflyBE

Go backend foundation for a flight price tracking app. This scaffold intentionally does not connect to Duffel yet and does not implement alerting yet.

## Requirements

- Go 1.26.1
- mockgen v0.6.0 installed on PATH

## Setup

```powershell
Copy-Item .env.example .env
```

The server loads `.env` directly for local development.

```powershell
go run ./cmd/server
```

The API defaults to `http://localhost:3000`.

## Development

Mocks are generated with mockgen into the repository and service package `mocks` directories.

```powershell
go install go.uber.org/mock/mockgen@v0.6.0
go generate ./internal/tracker/...
go test ./...
```

## Infrastructure

AWS infrastructure lives in `infra/` as a Go CDK app. See `infra/README.md` for synth and deploy commands.

## API

### Create a tracker

```http
POST /tracked-flights
Content-Type: application/json
```

```json
{
  "name": "Chicago to London",
  "user": {
    "email": "person@example.com"
  },
  "slices": [
    {
      "origin": "ORD",
      "destination": "LHR",
      "departure_date": "2026-09-10"
    },
    {
      "origin": "LHR",
      "destination": "ORD",
      "departure_date": "2026-09-18"
    }
  ],
  "passengers": [
    {
      "type": "adult"
    }
  ],
  "cabin_class": "economy"
}
```

### Read data

- `GET /tracked-flights`
- `GET /tracked-flights/{id}`
- `DELETE /tracked-flights/{id}`
