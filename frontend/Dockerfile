FROM node:20-alpine AS build
WORKDIR /app

COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
RUN apk add --no-cache curl
COPY --from=build /app/dist /usr/share/nginx/html
HEALTHCHECK --interval=5s --timeout=3s --retries=10 CMD curl -fsS http://localhost/ || exit 1
EXPOSE 80
