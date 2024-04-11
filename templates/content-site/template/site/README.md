# Site for {{ name }}

> {{ description }}

## Install 
1. Install [gvm](https://github.com/moovweb/gvm)
2. Run 
   ```bash
   gvm install go1.20.1
   gvm use go1.20.1
   ```
3. Go to `server` directory and run:
   ```bash
   make install 
   ```
4. Go to `client` directory and run:
    ```bash
    yarn
    ```

## Run dev server
1. Go to `server` directory and build server:
   ```bash
   make build
   ```
2. Run server from `server` directory:
   ```bash
   ./bin/linux_amd64/site
   ```
3. In another terminal go to `client` directory and run:
   ```bash
   yarn start
   ```
   
## Build production version
1. In `server` directory run:
   ```bash
   make build 
   ```

## Deploy
Create new linux user
```bash
sudo adduser {{ logoText }}-site
```
Copy `server/bin/linux_386/site` to `/home/{{ logoText }}-site/site/site` on server

Run `/home/{{ logoText }}-site/site/site` for create config.
Edit config `/home/{{ logoText }}-site/site/config.cfg`

Enable apache modules:
```bash
a2enmod proxy
a2enmod proxy_http
```

Edit apache config:
```
<VirtualHost 127.0.0.1:80>
    ServerName {{ logoText }}.com

    ProxyPreserveHost On
    ProxyRequests off
    ProxyPass / http://localhost:3002/
    ProxyPassReverse / http://localhost:3002/

    ErrorLog /var/log/apache2/error.{{ logoText }}.com.log
    CustomLog /var/log/apache2/access.{{ logoText }}.com.log combined
</VirtualHost>
```
Restart apache.

Install supervisor:
```bash
sudo apt-get -y install supervisor
```
Create log directory 
```bash
sudo mkdir -p /var/log/{{ logoText }}/site
```

Add in supervisor config `sudo nano /etc/supervisor/supervisord.conf`:
```
[program:{{ logoText }}-site]
directory=/home/{{ logoText }}-site/site/
command=/home/{{ logoText }}-site/site/site
stopasgroup=true
killasgroup=true
stopsignal=INT
autostart=true
autorestart=true
startsecs=10
stdout_logfile=/var/log/{{ logoText }}/site/stdout.log
stdout_logfile_maxbytes=1MB
stdout_logfile_backups=10
stdout_capture_maxbytes=1MB
stderr_logfile=/var/log/{{ logoText }}/site/stderr.log
stderr_logfile_maxbytes=1MB
stderr_logfile_backups=10
stderr_capture_maxbytes=1MB
environment = HOME="/home/{{ logoText }}-site", USER="{{ logoText }}-site"
user={{ logoText }}-site
```

Restart supervisor:
```bash
sudo service supervisor restart
```
