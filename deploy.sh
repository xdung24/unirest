#!/bin/sh
git pull &&  sh build.sh && echo 'binary built successfully' && \
sudo systemctl stop universal-rest.service && echo 'service stopped' && \
sudo cp ./tmp/main.exe /usr/local/bin/universal-rest && echo 'binary copied' && \
sudo systemctl restart universal-rest.service && echo 'service restarted' && \
sudo systemctl status universal-rest.service
