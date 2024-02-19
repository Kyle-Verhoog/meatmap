FROM golang:1.21.5-bookworm

WORKDIR /src

# Cache deps
ADD go.mod go.sum /src
RUN go mod download

ADD . /src
RUN go build -o meatmap
CMD ./meatmap
