FROM golang:1.13.9 AS builder

WORKDIR /workspace
COPY . .

RUN mkdir build
RUN CGO_ENABLED=0 go build -o ./build/client .


FROM scratch
WORKDIR /opt
COPY --from=builder /workspace/build/client .
COPY --from=builder /workspace/web ./web

ENTRYPOINT  [ "./client" ]
