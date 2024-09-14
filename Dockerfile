FROM golang:1.21.8 as icearrow
RUN mkdir -p /icearrow
WORKDIR /icearrow
COPY . .
RUN make

FROM jhaals/yopass:latest
COPY --from=icearrow /icearrow/walrus-client /icearrow/icearrow-server /
ENTRYPOINT ["/icearrow-server"]
