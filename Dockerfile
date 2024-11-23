FROM golang:1.23

WORKDIR /usr/src/app

COPY . .

# make life easy
CMD ["make run"]
