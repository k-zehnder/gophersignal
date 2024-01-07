# frontend

### Frontend (Update for Production):
```bash
# Install dependencies, build, and export the Next.js application
npm install && npm run build && npm run export

# Remove existing files in the web server's root directory
rm -rf /var/www/gophersignal.com/html/*

# Copy newly built static files to the web server's root directory
cp -r out/* /var/www/gophersignal.com/html/

# Restart the Nginx service to apply changes
sudo service nginx restart
```
