FROM golang:1.13.9 AS builder

WORKDIR /workspace
COPY . .

RUN mkdir build
RUN CGO_ENABLED=0 go build -o ./build/server_middleware_protected .


FROM scratch
WORKDIR /opt
COPY --from=builder /workspace/build/server_middleware_protected .

ENTRYPOINT  [ "./server_middleware_protected" ]
