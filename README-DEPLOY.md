## Location files

Config file: /etc/unirest.conf
Service file: /etc/systemd/system/unirest.service
Binary file: /usr/local/bin/unirest
nginx config: /etc/nginx/sites-enabled/unirest.conf
Your cert is in: /root/.acme.sh/lxd.ddns.net_ecc/lxd.ddns.net.cer
Your cert key is in: /root/.acme.sh/lxd.ddns.net_ecc/lxd.ddns.net.key
The intermediate CA cert is in: /root/.acme.sh/lxd.ddns.net_ecc/ca.cer
And the full chain certs is there: /root/.acme.sh/lxd.ddns.net_ecc/fullchain.cer


## Steps to deploy

- git pull
- go build
- copy binary to sudo mv /tmp/main.exe /usr/local/bin/unirest
- copy config file to sudo cp unirest.conf /etc/unirest/config.conf
- copy service file to sudo cp unirest-sample.conf /etc/systemd/system/unirest.service
- sudo systemctl daemon-reload
- sudo systemctl restart unirest.service
- update /etc/nginx/sites-enabled/unirest.conf
- sudo systemctl reload nginx
- sudo systemctl restart nginx.service


## steps to generate a new cert

acme.sh --install-cert -d lxd.ddns.net \
--key-file       /root/.acme.sh/lxd.ddns.net_ecc/lxd.ddns.net.key  \
--fullchain-file /root/.acme.sh/lxd.ddns.net_ecc/lxd.ddns.net.fullchain.cer \
--reloadcmd     "service nginx force-reload"