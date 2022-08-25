FROM golang:1.18

COPY dyna /usr/src/dyna
COPY config.json /usr/src/dyna/config.json

WORKDIR /usr/src/dyna

RUN ["make", "build_docker"]
ENTRYPOINT ["make", "r_entrypoint"]
