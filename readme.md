# Bookings App use Golang
# This is a simple bookings app using Golang
Use [chi](https://github.com/go-chi/chi) for routing
Use [nosurf](https://github.com/justinas/nosurf) for CSRF protection
Use [scs](https://github.com/alexedwards/scs/v2) for session management

## Installation
```bash
go get github.com/justinas/nosurf
go get github.com/go-chi/chi
go get github.com/alexedwards/scs/v2
```
## Or use go mod
```bash
go mod init
go mod tidy
```

## Run
```bash
go run cmd/web/*.go
```

## Build
```bash
go build -o bookings cmd/web/*.go
```