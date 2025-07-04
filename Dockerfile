FROM krakend/builder:2.9.4 AS builder

WORKDIR /app

COPY . .

RUN go build -buildmode=plugin -o krakend-grpc-proxy.so .

FROM krakend:2.9.4 AS runner

COPY --from=builder /app/*.so /etc/krakend-plugin