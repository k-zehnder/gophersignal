# Stage 1: Build the Next.js/React application
FROM node:latest as build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Stage 2: Setup Nginx to serve the static files and configure it
FROM nginx:latest
RUN rm -rf /usr/share/nginx/html/*
COPY ./nginx.conf /etc/nginx/nginx.conf
COPY --from=build /app/out /usr/share/nginx/html

# Expose ports for HTTP and HTTPS
EXPOSE 80
EXPOSE 443

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
