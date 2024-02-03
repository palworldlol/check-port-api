FROM golang:1.21 as build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /check-port-api

FROM golang:alpine

WORKDIR /

COPY --from=build-stage /check-port-api /check-port-api

EXPOSE 8181

# Run
ENTRYPOINT ["/check-port-api"]