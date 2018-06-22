# Builder stage for the project
FROM golang AS builder

ADD . /go/src/github.com/syntaqx/mixpanel
WORKDIR /go/src/github.com/syntaqx/mixpanel

RUN make dep
RUN make build

# Build a scratch container for the binary
FROM scratch
COPY --from=builder /go/src/github.com/syntaqx/mixpanel/bin/mixpanel /bin/mixpanel
CMD ["/bin/mixpanel"]
