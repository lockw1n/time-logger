# Playwright base image with browsers preinstalled
FROM mcr.microsoft.com/playwright:v1.57.0-jammy

WORKDIR /app

# Install deps
COPY package.json package-lock.json* ./
RUN npm ci --omit=dev || npm install --omit=dev

# Copy source
COPY server.mjs ./

ENV NODE_ENV=production
ENV PORT=3001

EXPOSE 3001

CMD ["node", "server.mjs"]
