FROM golang:1.21.8 as hippos
RUN mkdir -p /hippos
WORKDIR /hippos
COPY . .
RUN make

FROM jhaals/yopass:latest
COPY --from=hippos /hippos/walrus-client /hippos/hippo-server /
ENTRYPOINT ["/hippo-server"]
