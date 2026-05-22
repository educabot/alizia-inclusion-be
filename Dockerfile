FROM golang:1.26.1-alpine AS builder
RUN apk --no-cache add git
ARG GITHUB_TOKEN
RUN echo "machine github.com login x-token password ${GITHUB_TOKEN}" > /root/.netrc
ENV GOPRIVATE=github.com/educabot/*
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /alizia-inclusion-api ./cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates postgresql-client
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /alizia-inclusion-api .
COPY db/migrations ./db/migrations
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh
EXPOSE 8080
USER appuser
CMD ["./entrypoint.sh"]
