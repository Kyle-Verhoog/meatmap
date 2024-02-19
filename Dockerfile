FROM golang:1.21.5-bookworm

WORKDIR /src
ADD . /src
RUN go build -o meatmap
CMD ./meatmap
