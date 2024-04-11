# {{ name }}

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
   ./bin/linux_amd64/{{ shortName }}
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
sudo adduser {{ shortName }}
```
Copy `server/bin/linux_amd64/app` to `/home/{{ shortName }}/{{ name }}/{{ shortName }}` on server

Create DB & DB user:
```bash
sudo -u postgres createuser {{ shortName }}
sudo -u postgres createdb {{ shortName }}
```
Set permission & password
```
sudo -u postgres psql
psql=# alter user {{ shortName }} with encrypted password 'password';
psql=# grant all privileges on database {{ shortName }} to {{ shortName }};
```

Run `/home/{{ shortName }}/{{ name }}/{{ shortName }}` for create config.
Edit config `/home/{{ shortName }}/{{ name }}/config.cfg`

Enable apache modules:
```bash
a2enmod proxy
a2enmod proxy_http
a2enmod proxy_wstunnel
```

Edit apache config:
```
<VirtualHost 127.0.0.1:80>
    ServerName {{ shortName }}.servername.com

    ProxyPreserveHost On
    ProxyRequests off
    ProxyPass / http://localhost:3000/
    ProxyPassReverse / http://localhost:3000/

    <Location "/ws">
        ProxyPass "ws://localhost:3000/ws"
    </Location>

    ErrorLog /var/log/apache2/error.{{ name }}.log
    CustomLog /var/log/apache2/access.{{ name }}.log combined
</VirtualHost>
```
Restart apache.

Install supervisor:
```bash
sudo apt-get -y install supervisor
```
Create log directory 
```bash
sudo mkdir -p /var/log/{{ name }}
```

Add in supervisor config `sudo nano /etc/supervisor/supervisord.conf`:
```
[program:{{ shortName }}]
directory=/home/{{ shortName }}/{{ name }}/
command=/home/{{ shortName }}/{{ name }}/{{ shortName }}
autostart=true
autorestart=true
startsecs=10
stdout_logfile=/var/log/{{ name }}/stdout.log
stdout_logfile_maxbytes=1MB
stdout_logfile_backups=10
stdout_capture_maxbytes=1MB
stderr_logfile=/var/log/{{ name }}/stderr.log
stderr_logfile_maxbytes=1MB
stderr_logfile_backups=10
stderr_capture_maxbytes=1MB
environment = HOME="/home/{{ shortName }}", USER="{{ shortName }}"
user = {{ shortName }}
```

Restart supervisor:
```bash
sudo service supervisor restart
```
