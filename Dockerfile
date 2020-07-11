FROM golang:1.14 AS builder
ENV GOPROXY="https://goproxy.cn,direct"
ENV GO111MODULE=on
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/smzdm-reptile

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=builder /src/smzdm-reptile /smzdm-reptile
CMD ["/smzdm-reptile"]