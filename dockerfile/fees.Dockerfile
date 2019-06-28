FROM golang:1.11.4-alpine as builder

RUN apk add --no-cache git dep openssh-client

WORKDIR /go/src/github.com/Ankr-network/dccn-fees
COPY . .

RUN dep ensure -vendor-only


RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cmd/fees main.go

RUN echo '0 2 * * * go run /go/src/github.com/Ankr-network/dccn-fees/crontab/daily-fees.go  >>/var/log/daily.log & ' >> /etc/crontabs/root
RUN echo '0 3 1 * * go run /go/src/github.com/Ankr-network/dccn-fees/crontab/monthly-fees.go  >>/var/log/monthly.log & ' >> /etc/crontabs/root
RUN echo '0 4 1 * * go run /go/src/github.com/Ankr-network/dccn-fees/crontab/monthly-clear.go  >>/var/log/monthly_clear.log & ' >> /etc/crontabs/root

COPY dockerfile/wrapper_script.sh wrapper_script.sh

CMD ./wrapper_script.sh


