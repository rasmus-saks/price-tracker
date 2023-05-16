FROM golang:1.20 as builder

ENV PROJECT github.com/rasmus-saks/price-tracker
ENV GO111MODULE on
WORKDIR /go/src/$PROJECT

COPY go.mod /go/src/$PROJECT
COPY go.sum /go/src/$PROJECT

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

FROM alpine
COPY --from=builder /main /main

ENTRYPOINT ["/main"]
