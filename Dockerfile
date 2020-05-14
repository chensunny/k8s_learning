FROM golang:1.12.13-buster

RUN mkdir -p /deploy

ENV WORK_DIR /deploy/web
WORKDIR  ${WORK_DIR}

RUN mkdir -p ${WORK_DIR}

COPY . .
RUN sh ./scripts/build-bin.sh

ENV BINARY=web
EXPOSE 8080

CMD ["sh", "-c", "./${BINARY}"]
