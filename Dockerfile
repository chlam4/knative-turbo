# Set the base image

FROM alpine:3.3

# Set the file maintainer

MAINTAINER Pallavi Debnath <pallavi.debnath@turbonomic.com>


ADD _output/knative-turbo.linux /bin/knativeturbo


ENTRYPOINT ["/bin/knativeturbo"]
