FROM node:12-buster as frontend

WORKDIR /app
COPY web .
RUN npm install --quiet && npm run build

FROM golang:1.15-alpine AS backend

RUN apk add --no-cache ca-certificates git
WORKDIR /go/src/github.com/geekgonecrazy/prismplus/
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build

FROM alpine

RUN mkdir /app
WORKDIR /app

COPY --from=frontend /app /app/web

COPY --from=backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend /go/src/github.com/geekgonecrazy/prismplus/prismplus /app/prismplus

EXPOSE 5383
EXPOSE 1935

CMD ["/app/prismplus"]