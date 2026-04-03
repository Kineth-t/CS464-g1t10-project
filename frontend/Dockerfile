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

# Copy nginx config with env variables to templates
COPY nginx.conf /etc/nginx/templates/default.conf.template

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]