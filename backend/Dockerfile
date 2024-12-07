FROM docker.io/library/golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o KAM_backend -ldflags="-s -w" .

FROM scratch AS runtime

WORKDIR /app
COPY --from=build /app/KAM_backend /app/
USER 1000
EXPOSE 8080
EXPOSE 8082
ENV IN_CLUSTER_MODE=false
ENTRYPOINT ["/app/KAM_backend"]
CMD ["$(if [ \"$IN_CLUSTER_MODE\" = \"true\" ]; then echo \"-in-cluster\"; fi)"]
