
FROM golang:alpine AS builder

RUN apk add --update alpine-sdk ca-certificates && update-ca-certificates

WORKDIR /opt/build
COPY ./ ./

RUN go mod download

# With CGO:
# RUN GOOS=linux GOARCH=amd64 go build -tags musl -ldflags '-linkmode external -w -s -extldflags "-static"' -o /opt/build/coding-exercise cmd/exercise/main.go

# Without CGO:
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o /opt/build/coding-exercise cmd/exercise/main.go

RUN mkdir -p /out
RUN mkdir -p /in

ENV USER=appuser
ENV UID=10001

RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

FROM scratch 

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /opt/build/coding-exercise /go/bin/coding-exercise

COPY --from=builder /out /out
COPY --from=builder /in /in

USER appuser:appuser

# Set max cpu to 1 so you can control the cpu usage and make the go service aware it is limited
ENV GOMAXPROCS=1

ENTRYPOINT ["/go/bin/coding-exercise"]
