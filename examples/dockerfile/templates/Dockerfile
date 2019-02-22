FROM golang:1.11-alpine3.9 as builder

WORKDIR /project

{{if len .addons}}
RUN set -xe && \
    apk update && apk upgrade && \
    apk add --no-cache make {{if has "certs" .addons}}ca-certificates {{end}}git curl{{if has "migrate" .addons}} && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v3.5.2/migrate.linux-amd64.tar.gz | tar xvz && \
    cp migrate.linux-amd64 /migrate{{end}}
{{end}}

COPY . .

RUN make dep && \
    make build-{{.svc}}{{if has "health" .addons}} && \
    make build-health{{end}}

FROM scratch

{{if has "certs" .addons}}COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/{{end}}
{{if has "migrate" .addons}}COPY --from=builder /migrate /migrate
COPY --from=builder /project/internal/{{.svc}}/migrations /migrations{{end}}
COPY --from=builder /project/dist/{{.svc}} /{{.svc}}
{{if has "assets" .addons}}COPY --from=builder /project/internal/notify/assets /assets{{end}}
{{if has "health" .addons}}COPY --from=builder /project/dist/health /health

HEALTHCHECK CMD ["/health"]{{end}}

{{if len .expose}}EXPOSE{{if has "metrics" .expose}} 8080{{end}}{{if has "grpc" .expose}} 8086{{end}}{{end}}

CMD ["/{{.svc}}{{if .cmd}} {{.cmd}}{{end}}"]
