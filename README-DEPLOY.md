## Location files

Config file: /etc/universal-rest.conf
Service file: /etc/systemd/system/universal-rest.service
Binary file: /usr/local/bin/universal-rest
nginx config: /etc/nginx/sites-enabled/universal-rest.conf
Your cert is in: /root/.acme.sh/lxd.ddns.net_ecc/lxd.ddns.net.cer
Your cert key is in: /root/.acme.sh/lxd.ddns.net_ecc/lxd.ddns.net.key
The intermediate CA cert is in: /root/.acme.sh/lxd.ddns.net_ecc/ca.cer
And the full chain certs is there: /root/.acme.sh/lxd.ddns.net_ecc/fullchain.cer


## Steps to deploy

- git pull
- go build
- copy binary to sudo mv /tmp/main.exe /usr/local/bin/universal-rest
- copy config file to sudo cp universal-rest.conf /etc/universal-rest/config.conf
- copy service file to sudo cp universal-rest-sample.conf /etc/systemd/system/universal-rest.service
- sudo systemctl daemon-reload
- sudo systemctl restart universal-rest.service
- update /etc/nginx/sites-enabled/universal-rest.conf
- sudo systemctl reload nginx
- sudo systemctl restart nginx.service


## steps to generate a new cert

acme.sh --install-cert -d lxd.ddns.net \
--key-file       /root/.acme.sh/lxd.ddns.net_ecc/lxd.ddns.net.key  \
--fullchain-file /root/.acme.sh/lxd.ddns.net_ecc/lxd.ddns.net.fullchain.cer \
--reloadcmd     "service nginx force-reload"