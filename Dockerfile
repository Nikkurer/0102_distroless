# Start by building the application.
FROM golang:1.19-alpine as builder

WORKDIR /go/src/app
COPY . .

RUN echo $DB_URL
RUN go mod download
RUN go build -o /go/bin/app cmd/main.go

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian11

ENV HOST=""
ENV PORT=""
ENV DB_URL=""
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /lib/ld-musl-aarch64.so.1 /lib/ld-musl-aarch64.so.1
USER nobody
COPY --from=builder --chown=nobody /go/bin/app /bin/app
ENTRYPOINT ["/bin/app", "-port=${PORT}", "-host=${HOST}", "-dbUrl=$DB_URL"]
