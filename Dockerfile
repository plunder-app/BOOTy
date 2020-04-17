FROM golang:1.13-alpine as dev
RUN apk add --no-cache git ca-certificates gcc linux-headers musl-dev
RUN adduser -D appuser
COPY . /src/
WORKDIR /src
RUN go get github.com/schollz/progressbar 
RUN go get github.com/micmonay/keybd_event
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o booty

FROM scratch
COPY --from=dev /src/booty /
ENTRYPOINT ["/booty"]