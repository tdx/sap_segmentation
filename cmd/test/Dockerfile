FROM golang:alpine AS builder

WORKDIR /go/src/sap_segmentation

ENV GOPATH /go

COPY go.mod .
COPY go.sum .
RUN go mod download

ADD . .

RUN cd cmd/sap_segmentation && go build -o /go/bin/sap_segmentation .

FROM alpine:latest
COPY --from=builder /go/bin/sap_segmentation /bin/sap_segmentation

ENTRYPOINT ["/bin/sap_segmentation"]
