FROM alpine:latest
RUN apk update && apk add netcat-openbsd
COPY /netcat-client /
RUN chmod u+x test-server-response.sh
ENTRYPOINT [ "./test-server-response.sh" ]
