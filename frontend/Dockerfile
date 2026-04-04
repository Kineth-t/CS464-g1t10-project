# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy frontend files
COPY package*.json ./
RUN npm install

COPY . .

# Accept VITE_API_URL as a build argument so Vite can embed it at build time
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

# Build for production
RUN npm run build

# Run stage - serve with nginx
FROM nginx:alpine

# Copy built files to nginx html directory
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx config template — BACKEND_URL is substituted at container startup
COPY nginx.conf.template /etc/nginx/conf.d/default.conf.template

EXPOSE 8080

# Substitute env vars into nginx config, then start nginx
CMD ["/bin/sh", "-c", "envsubst '${BACKEND_URL}' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"]