# syntax=docker/dockerfile:1

FROM golang:1.24.3 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/bokwoon95/wgo@latest
COPY . ./

CMD [ "wgo", "run", "main.go" ]
# RUN CGO_ENABLED=0 GOOS=linux go build -o /execstasy-api

# Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

# # Deploy the application binary into a lean image
# FROM gcr.io/distroless/base-debian11 AS build-release-stage

# WORKDIR /

# COPY --from=build-stage /execstasy-api /execstasy-api

# EXPOSE 5000

# USER nonroot:nonroot

# CMD ["/execstasy-api"]

