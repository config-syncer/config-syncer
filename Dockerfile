FROM debian:jessie-20190506
COPY kubed /usr/bin/kubed
ENTRYPOINT ["/usr/bin/kubed"]

