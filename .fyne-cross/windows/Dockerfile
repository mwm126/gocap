FROM fyneio/fyne-cross:1.1-windows-22.01.10
LABEL maintainer="Mark Meredith <mark.meredith@pm.me>"

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get -qq update \
    && apt-get -qq -y install --no-install-recommends \
    libusb-1.0=2:1.* \
    libyubikey-dev=1.* \
    libykpers-1-dev=1.* \
    && rm -rf /var/lib/apt \
    && curl --insecure "https://developers.yubico.com/yubikey-personalization/Releases/ykpers-1.20.0-win64.zip" --output ykpers.zip \
    && unzip -d /usr/x86_64-w64-mingw32 ykpers.zip
