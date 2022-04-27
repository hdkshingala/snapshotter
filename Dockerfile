FROM golang:1.18.0
WORKDIR /snapshotter
ADD . .
RUN go mod download && CGO_ENABLED=0 go build

FROM scratch
WORKDIR /snapshotter
COPY --from=0 snapshotter .
ENTRYPOINT [ "./snapshotter" ]
