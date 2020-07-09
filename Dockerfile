FROM golang:1.14.4-alpine3.12 as builder

# install depencencies.
RUN apk add --update gcc musl-dev

ENV APPDIR=github.com/itomofumi/go-gin-xorm-starter

WORKDIR /go/src/${APPDIR}
RUN apk add --update git make ca-certificates gcc musl-dev openssh-client

# Setup SSH
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"
RUN mkdir -p -m 600 ~/.ssh \
    && ssh-keyscan -H github.com >> ~/.ssh/known_hosts

COPY . .
# RUN make lint
# RUN make test
RUN make build
RUN cp ./bin/starter /tmp/starter

FROM alpine

COPY --from=builder /tmp/starter /opt/api/starter
RUN apk add --no-cache ca-certificates
ENV GIN_MODE=release
ENV PORT=80
CMD ["/opt/api/starter"]
