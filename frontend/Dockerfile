# frontend/Dockerfile

# Stage 1: Build React/Vite
FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

# --- BAGIAN PENTING ---
# 1. Terima argumen dari docker-compose
ARG VITE_API_URL
ARG VITE_APP_TITLE

# 2. Set sebagai Environment Variable agar Vite bisa membacanya saat build
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_APP_TITLE=$VITE_APP_TITLE
RUN echo "URL API adalah: $VITE_API_URL"
# ----------------------

RUN npm run build

# Stage 2: Serve dengan Nginx
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
# Catatan: Vite biasanya output ke folder /dist (bukan /build)

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]