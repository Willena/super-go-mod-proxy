FROM golang:1.16-alpine as build
WORKDIR /build/
ADD . .
RUN go mod download
RUN go build -o /tmp/super-go-mod-proxy github.com/willena/super-go-mod-proxy/cmd


FROM alpine
WORKDIR /app/
COPY --from=build /tmp/super-go-mod-proxy .
RUN chmod +x ./super-go-mod-proxy

ENTRYPOINT ["/app/super-go-mod-proxy"]

