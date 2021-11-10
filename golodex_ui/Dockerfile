# syntax=docker/dockerfile:1

FROM node:14.17.6-alpine AS builder

WORKDIR /app

COPY package.json ./

RUN npm install

COPY package-lock.json ./
COPY public/ ./public/
COPY src/ ./src/
COPY .env_production ./.env

RUN npm run build

#Stage 2
#######################################
#pull the official nginx:1.19.0 base image
FROM nginx:1.19.0
#copies React to the container directory
# Set working directory to nginx resources directory
WORKDIR /usr/share/nginx/html
# Remove default nginx static resources
RUN rm -rf ./*
# Copies static resources from builder stage
COPY --from=builder /app/build .
# needed this to make React Router work properly
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx/nginx.conf /etc/nginx/conf.d

EXPOSE 80
# Containers run nginx with global directives and daemon off
ENTRYPOINT ["nginx", "-g", "daemon off;"]