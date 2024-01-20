# Development

## Setup
Create .env based on .env.template

## Commands

### Create empty migration
```bash
migrate create -ext sql -dir schemata <migration_name>
```

### Run migration
```bash
go run cmd/migrate/main.go
```

### Create models
```bash
sqlc generate
```