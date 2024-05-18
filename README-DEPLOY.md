## Location files

Config file: /etc/universal-rest.conf
Service file: /etc/systemd/system/universal-rest.service
Binary file: /usr/local/bin/universal-rest
Your cert is in: /root/.acme.sh/lxd.ddns.net_ecc/lxd.ddns.net.cer
Your cert key is in: /root/.acme.sh/lxd.ddns.net_ecc/lxd.ddns.net.key
The intermediate CA cert is in: /root/.acme.sh/lxd.ddns.net_ecc/ca.cer
And the full chain certs is there: /root/.acme.sh/lxd.ddns.net_ecc/fullchain.cer


## Steps to deploy

- git pull
- go build
- copy binary to /usr/local/bin/universal-rest
- copy config file to /etc/universal-rest.conf
- copy service file to /etc/systemd/system/universal-rest.service
- sudo systemctl daemon-reload
- sudo systemctl start universal-rest.service
