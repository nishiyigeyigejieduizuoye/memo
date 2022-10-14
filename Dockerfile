FROM golang:1.18-alpine AS build
RUN apk add --no-cache gcc upx
COPY ./ /app/
WORKDIR /app/
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download && go build -o main && upx main

FROM node:current-alpine AS frontend-build
COPY frontend/ /app/
WORKDIR /app/
RUN yarn install && yarn build

FROM alpine:latest
WORKDIR /app
VOLUME /app/data
COPY --from=frontend-build /app/dist /app/frontend
COPY --from=build /app/main /app/
ENV GIN_MODE=release
ENTRYPOINT /app/main
EXPOSE 80/tcp
