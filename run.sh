#!/bin/bash
set -e

APPENV=${APPENV:-stinkbaitenv}

/opt/bin/s3kms -r us-west-1 get -b opsee-keys -o dev/$APPENV > /$APPENV

source /$APPENV && \
  /opt/bin/s3kms -r us-west-1 get -b opsee-keys -o dev/$STINKBAIT_CERT > /$STINKBAIT_CERT && \
  /opt/bin/s3kms -r us-west-1 get -b opsee-keys -o dev/$STINKBAIT_CERT_KEY > /$STINKBAIT_CERT_KEY && \
  chmod 600 /$STINKBAIT_CERT_KEY && \
  /stinkbait
