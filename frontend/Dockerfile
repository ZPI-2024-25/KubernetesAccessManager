FROM node as vite-app

WORKDIR /app/
COPY . .

RUN ["npm", "i"]
RUN ["npm", "run", "build"]

FROM nginx:alpine

RUN apk add --no-cache gettext

COPY --from=vite-app /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/templates/nginx.conf.template

ENV BACKEND_SERVICE_HOST=kam-backend
ENV BACKEND_SERVICE_PORT=8080

EXPOSE 80

CMD ["/bin/sh", "-c", "envsubst '$BACKEND_SERVICE_HOST $BACKEND_SERVICE_PORT' < /etc/nginx/templates/nginx.conf.template > /etc/nginx/conf.d/default.conf && exec nginx -g 'daemon off;'"]
