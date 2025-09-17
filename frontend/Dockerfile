# build environment
FROM node:20-alpine AS build

ARG VERSION
ARG GIT_COMMIT

WORKDIR /app
ENV PATH=/app/node_modules/.bin:$PATH
COPY package.json ./
COPY package-lock.json ./
RUN npm config set fetch-retry-mintimeout 100000
RUN npm config set fetch-retry-maxtimeout 600000
RUN npm ci --no-audit
COPY . ./

RUN echo "VITE_FRONTEND_VERSION=$VERSION" >> .env.production
RUN echo "VITE_FRONTEND_BUILD_TIME=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" >> .env.production
RUN echo "VITE_FRONTEND_GIT_COMMIT=$GIT_COMMIT" >> .env.production

RUN npm run build

# production environment
FROM nginx:stable-alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY --from=build /app/nginx/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
