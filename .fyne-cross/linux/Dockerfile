FROM fyneio/fyne-cross:1.3-base-22.07.13

LABEL maintainer="Mark Meredith <mark.meredith@pm.me>"

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get -qq update \
    && apt-get -qq -y install --no-install-recommends \
    libusb-1.0=2:1.* \
    libyubikey-dev=1.* \
    libykpers-1-dev=1.* \
    && rm -rf /var/lib/apt
