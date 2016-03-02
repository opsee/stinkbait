FROM alpine:3.3

ENV STINKBAIT_ADDRESS ":9100"
ENV STINKBAIT_CERT "cert.pem"
ENV STINKBAIT_CERT_KEY "key.pem"
ENV STINKBAIT_MEMCACHED_NODES "memcached:11211"
ENV APPENV "testenv"

RUN apk add --update bash ca-certificates curl && \
	rm -rf /var/cache/apk/*

RUN mkdir -p /opt/bin && \
		curl -Lo /opt/bin/s3kms https://s3-us-west-2.amazonaws.com/opsee-releases/go/vinz-clortho/s3kms-linux-amd64 && \
    chmod 755 /opt/bin/s3kms

COPY run.sh /
COPY target/linux/amd64/bin/* /
COPY cert.pem /
COPY key.pem /

EXPOSE 9100
CMD ["/stinkbait"]
