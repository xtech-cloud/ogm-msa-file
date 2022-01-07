# *************************************
#
# OpenGM
#
# *************************************

FROM alpine:3.14

MAINTAINER XTech Cloud "xtech.cloud"

ENV container docker
ENV MSA_MODE release

EXPOSE 18808

ADD bin/ogm-file /usr/local/bin/
RUN chmod +x /usr/local/bin/ogm-file

CMD ["/usr/local/bin/ogm-file"]
