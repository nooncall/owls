FROM centos:latest

LABEL MAINTAINER=developer@owls.nooncall.cn

COPY ./bin /service/

WORKDIR /service

ENTRYPOINT ["./owls"]
