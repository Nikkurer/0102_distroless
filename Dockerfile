# Start by building the application.
FROM golang:1.19-alpine as builder

WORKDIR /go/src/app
COPY . .

RUN echo $DB_URL
RUN go mod download
RUN go build -o /go/bin/app cmd/main.go

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian11
ARG HOST
ENV HOST=$HOST
ARG PORT
ENV PORT=$PORT
ARG DB_URL
ENV DB_URL=$DB_URL
EXPOSE 9000
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /bin/sh /bin/sh
COPY --from=builder /lib/ld-musl-aarch64.so.1 /lib/ld-musl-aarch64.so.1
USER nobody
COPY --from=builder --chown=nobody /go/bin/app /bin/app
ENTRYPOINT ["/bin/sh", "-c", "/bin/app", "-port=$PORT", "-host=$HOST", "-dbUrl=$DB_URL"]
