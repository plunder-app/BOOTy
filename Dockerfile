FROM golang:1.13-alpine as dev
RUN apk add --no-cache git ca-certificates gcc linux-headers musl-dev
RUN adduser -D appuser
COPY . /go/src/github.com/thebsdbox/BOOTy/
WORKDIR /go/src/github.com/thebsdbox/BOOTy
RUN go get github.com/schollz/progressbar 
RUN go get github.com/micmonay/keybd_event
RUN go get github.com/spf13/cobra
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o booty

FROM scratch
COPY --from=dev /go/src/github.com/thebsdbox/BOOTy/booty /
ENTRYPOINT ["/booty"]