# Start by building the application.
FROM golang:1.19-alpine as builder

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go build -o /go/bin/app cmd/main.go

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian11
EXPOSE $PORT
ARG ARCH=$(uname -m)
COPY --from=builder /etc/passwd /etc/group /etc/
COPY --from=builder /bin/sh /bin/sh
COPY --from=builder /lib/ld-musl-* /lib/
USER nobody:nobody
COPY --from=builder --chown=nobody:nobody /go/bin/app /bin/app
ENTRYPOINT ["/bin/sh", "-c", "/bin/app -port=${PORT} -host=${HOST} -dbUrl=${DB_URL}"]
