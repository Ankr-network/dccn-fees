FROM golang:1.13-alpine as builder

# privaite repo go module
RUN apk add --no-cache git dep openssh-client
ARG	GITHUB_USER
ARG	GITHUB_TOKEN
RUN echo "machine github.com login ${GITHUB_USER} password ${GITHUB_TOKEN}" > ~/.netrc

WORKDIR /go/src/github.com/Ankr-network/dccn-fees
COPY . .

RUN dep ensure -vendor-only


RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cmd/fees main.go

RUN echo '0 2 * * * go run /go/src/github.com/Ankr-network/dccn-fees/crontab/daily-fees.go  >>/var/log/daily.log & ' >> /etc/crontabs/root
RUN echo '0 3 1 * * go run /go/src/github.com/Ankr-network/dccn-fees/crontab/monthly-fees.go  >>/var/log/monthly.log & ' >> /etc/crontabs/root
RUN echo '0 4 1 * * go run /go/src/github.com/Ankr-network/dccn-fees/crontab/monthly-clear.go  >>/var/log/monthly_clear.log & ' >> /etc/crontabs/root
RUN echo '1 3 1 * * go run /go/src/github.com/Ankr-network/dccn-fees/crontab/monthly-fees-for-provider.go  >>/var/log/monthly_for_provider.log & ' >> /etc/crontabs/root
RUN echo '1 4 1 * * go run /go/src/github.com/Ankr-network/dccn-fees/crontab/monthly-clear-for-provider.go  >>/var/log/monthly_clear_for_provider.log & ' >> /etc/crontabs/root


COPY dockerfile/wrapper_script.sh wrapper_script.sh

CMD ./wrapper_script.sh


