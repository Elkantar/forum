#dockerfile ascii-art-web
# sudo chmod 666 /var/run/docker.sock   
# docker run --publish 8050:8050
FROM golang:latest
LABEL Authors=Gregory/Thibault/Kenny/Maxence/Mathis
RUN mkdir /app
ADD . /app
WORKDIR /app
# copy all the document
COPY . .
RUN go build -o forum .
CMD ["./forum"]
EXPOSE 8050"
