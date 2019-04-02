FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
MAINTAINER Lucas Servén Marín <lserven@gmail.com>
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY bin/drae /
ENTRYPOINT ["/drae"]
CMD ["api"]
