winston@vps-ca69b6ab:~/code/gophersignal/backend$ cat /etc/nginx/sites-available/gophersignal.com 
server {

        root /var/www/gophersignal.com/html;
        index index.html index.htm index.nginx-debian.html;

        server_name gophersignal.com www.gophersignal.com;

        location / {
                try_files $uri $uri/ =404;
        }

	    location /articles {
		proxy_pass http://localhost:8080;
		proxy_http_version 1.1;
		proxy_set_header Upgrade $http_upgrade;
		proxy_set_header Connection 'upgrade';
		proxy_set_header Host $host;
		proxy_cache_bypass $http_upgrade;
	    }

    listen [::]:443 ssl ipv6only=on; # managed by Certbot
    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/gophersignal.com/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/gophersignal.com/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot


}
server {
    if ($host = www.gophersignal.com) {
        return 301 https://$host$request_uri;
    } # managed by Certbot


    if ($host = gophersignal.com) {
        return 301 https://$host$request_uri;
    } # managed by Certbot


        listen 80;
        listen [::]:80;

        server_name gophersignal.com www.gophersignal.com;
    return 404; # managed by Certbot




}