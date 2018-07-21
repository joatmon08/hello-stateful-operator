FROM alpine:3.6

RUN adduser -D hello-stateful-operator
USER hello-stateful-operator

ADD tmp/_output/bin/hello-stateful-operator /usr/local/bin/hello-stateful-operator
