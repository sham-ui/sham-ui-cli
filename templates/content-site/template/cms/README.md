# {{ name }}

> {{ description }}

## Install 
1. Install [gvm](https://github.com/moovweb/gvm)
2. Run 
   ```bash
   gvm install go1.20.2
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
   ./bin/linux_amd64/cms
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
sudo adduser {{ logoText }}-cms
```
Copy `server/bin/linux_386/cms` to `/home/{{ logoText }}-cms/cms/cms` on server

Create DB & DB user:
```bash
sudo -u postgres createuser {{ dbUser }}
sudo -u postgres createdb {{ dbName }}
```
Set permission & password
```
sudo -u postgres psql
psql=# alter user {{ dbUser }} with encrypted password 'password';
psql=# grant all privileges on database {{ dbUser }} to {{ dbUser }};
```

Run `/home/{{ logoText }}-cms/cms/cms` for create config.
Edit config `/home/{{ logoText }}-cms/cms/config.cfg`

Enable apache modules:
```bash
a2enmod proxy
a2enmod proxy_http
```

Edit apache config:
```
<VirtualHost 127.0.0.1:80>
    ServerName {{ logoText }}-cms.servername.com

    ProxyPreserveHost On
    ProxyRequests off
    ProxyPass / http://localhost:3001/
    ProxyPassReverse / http://localhost:3001/

    ErrorLog /var/log/apache2/error.{{ logoText }}-cms.servername.com.log
    CustomLog /var/log/apache2/access.{{ logoText }}-cms.servername.com.log combined
</VirtualHost>
```
Restart apache.

Install supervisor:
```bash
sudo apt-get -y install supervisor
```
Create log directory 
```bash
sudo mkdir -p /var/log/{{ logoText }}/cms
```

Add in supervisor config `sudo nano /etc/supervisor/supervisord.conf`:
```
[program:{{ logoText }}-cms]
directory=/home/{{ logoText }}-cms/cms/
command=/home/{{ logoText }}-cms/cms/cms
autostart=true
autorestart=true
startsecs=10
stdout_logfile=/var/log/{{ logoText }}/cms/stdout.log
stdout_logfile_maxbytes=1MB
stdout_logfile_backups=10
stdout_capture_maxbytes=1MB
stderr_logfile=/var/log/{{ logoText }}/cms/stderr.log
stderr_logfile_maxbytes=1MB
stderr_logfile_backups=10
stderr_capture_maxbytes=1MB
environment=HOME="/home/{{ logoText }}-cms", USER="{{ logoText }}-cms"
user={{ logoText }}-cms
umask=000
```

Restart supervisor:
```bash
sudo service supervisor restart
```
