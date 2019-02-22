FROM golang:1.11-alpine3.9 as builder

WORKDIR /project


RUN set -xe && \
    apk update && apk upgrade && \
    apk add --no-cache make ca-certificates git curl


COPY . .

RUN make dep && \
    make build-example

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /project/dist/example /example



EXPOSE 8080

CMD ["/example"]
