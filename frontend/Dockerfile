# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy frontend files
COPY package*.json ./
RUN npm install

COPY . .

# Build for production
RUN npm run build

# Run stage - serve with nginx
FROM nginx:alpine

# Copy built files to nginx html directory
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx config to templates for env substitution
COPY nginx.conf /etc/nginx/templates/default.conf.template

# Create startup script that substitutes environment variables
RUN echo '#!/bin/sh' > /docker-entrypoint.d/99-envsubst.sh && \
    echo 'export BACKEND_URL' >> /docker-entrypoint.d/99-envsubst.sh && \
    echo 'envsubst < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf' >> /docker-entrypoint.d/99-envsubst.sh && \
    chmod +x /docker-entrypoint.d/99-envsubst.sh

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]