# From alpine
FROM alpine:latest
# Install chrome
RUN echo @edge http://nl.alpinelinux.org/alpine/edge/community >> /etc/apk/repositories \
    && echo @edge http://nl.alpinelinux.org/alpine/edge/main >> /etc/apk/repositories \
    && apk add --no-cache \
    chromium@edge \
    harfbuzz@edge \
    nss@edge \
    freetype@edge \
    ttf-freefont@edge \
    libstdc++@edge \
    && rm -rf /var/cache/* \
    && mkdir /var/cache/apk
# chromium-browser --headless --no-sandbox --disable-gpu --print-to-pdf=/tmp/test.pdf https://www.chromestatus.com/

