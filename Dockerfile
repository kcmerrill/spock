FROM golang:1.7.5
MAINTAINER kc merrill <kcmerrill@gmail.com>
COPY . /go/src/github.com/kcmerrill/spock
WORKDIR /go/src/github.com/kcmerrill/spock
RUN  go build -ldflags "-X main.Commit=`git rev-parse HEAD` -X main.Version=0.1.`git rev-list --count HEAD`" -o /usr/local/bin/spock
RUN mkdir /spock/
EXPOSE 80
ENTRYPOINT ["spock"]
CMD ["--dir", "/spock/"]