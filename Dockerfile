FROM tetafro/golang-gcc:1.16-alpine AS build
RUN apk add --no-cache upx
COPY ./ /app/
WORKDIR /app/
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download && go build -o main && upx main

FROM alpine:latest
WORKDIR /app
VOLUME /app/data
COPY --from=build /app/main /app/
ENV GIN_MODE=release
ENTRYPOINT /app/main
EXPOSE 4000/tcp
