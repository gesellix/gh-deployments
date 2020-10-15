FROM alpine:3.12 AS builder
LABEL builder=true

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV APPPATH /app

RUN adduser -DH user
RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc ca-certificates
COPY . $APPPATH
RUN cd $APPPATH && go get -d \
 && go test -short ./... \
 && go build \
    -a \
    -ldflags '-s -w -extldflags "-static"' \
    -o /bin/gh-deployments

FROM scratch
LABEL maintainer="Tobias Gesellchen <tobias@gesellix.de> (@gesellix)"

ENTRYPOINT [ "/gh-deployments" ]
CMD [ ]

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
USER user

COPY --from=builder /bin/gh-deployments /gh-deployments
