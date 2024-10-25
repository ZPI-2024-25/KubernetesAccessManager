FROM node as vite-app

WORKDIR /app/
COPY . .

RUN ["npm", "i"]
RUN ["npm", "run", "build"]

FROM nginx:alpine

COPY --from=vite-app /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

ENTRYPOINT ["nginx", "-g", "daemon off;"]