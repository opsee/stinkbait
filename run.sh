#!/bin/bash
set -e

APPENV=${APPENV:-spanxenv}

/opt/bin/s3kms -r us-west-1 get -b opsee-keys -o dev/$APPENV > /$APPENV

/opt/bin/s3kms -r us-west-1 get -b opsee-keys -o dev/$APPENV > /$APPENV

source /$APPENV && \
  /opt/bin/s3kms -r us-west-1 get -b opsee-keys -o dev/$SPANX_CERT > /$SPANX_CERT && \
  /opt/bin/s3kms -r us-west-1 get -b opsee-keys -o dev/$SPANX_CERT_KEY > /$SPANX_CERT_KEY && \
  chmod 600 /$SPANX_CERT_KEY && \
	/opt/bin/migrate -url "$POSTGRES_CONN" -path /migrations up && \
  /spanx
