FROM node:22-alpine AS builder
WORKDIR /app
COPY ./frontend/package*.json ./frontend/tsconfig*.json ./
RUN npm ci
COPY ./frontend/ ./
RUN npm run build

FROM node:22-alpine
WORKDIR /app
COPY ./frontend/package*.json ./frontend/vite.config.ts ./
RUN npm ci --omit=dev
COPY --from=builder /app/build/ ./build/
ENTRYPOINT ["npm", "run", "preview"]