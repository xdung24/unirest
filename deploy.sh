#!/bin/sh
git pull &&  sh build.sh && echo 'binary built successfully' && \
sudo systemctl stop unirest.service && echo 'service stopped' && \
sudo cp ./tmp/main.exe /usr/local/bin/unirest && echo 'binary copied' && \
sudo systemctl restart unirest.service && echo 'service restarted' && \
sudo systemctl status unirest.service
