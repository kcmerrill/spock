# Spock

Making sure your applications live long and prosper.

[![Build Status](https://travis-ci.org/kcmerrill/spock.svg?branch=master)](https://travis-ci.org/kcmerrill/spock) [![Join the chat at https://gitter.im/kcmerrill/spock](https://badges.gitter.im/kcmerrill/spock.svg)](https://gitter.im/kcmerrill/spock?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

![Spock](assets/spock3.jpg "Spock")

## Binaries || Installation

[![MacOSX](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/apple_logo.png "Mac OSX")](http://go-dist.kcmerrill.com/kcmerrill/spock/mac/amd64) [![Linux](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/linux_logo.png "Linux")](http://go-dist.kcmerrill.com/kcmerrill/spock/linux/amd64)

via go:

`$ go get -u github.com/kcmerrill/spock`

via docker:

`$ docker run -ti -v $PWD/dir/to/root/channels/and/checks:/spock kcmerrill/spock`

## Checks and Channels

The concept behind `spock` makes it dead simple to get up and running. There are two aspects you should familiarize yourself with before continuing on. `checks` and `channels`. Essentially they are yaml file(s) located in their respective `checks` and `channels` folders.

### Checks

Checks are actions that `spock` does at regular intervals. You can use `cron` syntax, `every` syntax. You can [read more about it here](https://godoc.org/github.com/robfig/cron). What `spock` does at these intervals is completely up to you. Custom checks, url checks, disk space, etc ... it's up to you and what you need to check and alert on.

A quick example of what checks would look like. Within the `checks` folder you create, you can have one or many yaml files, just be sure to keep the check names unique. The checkname in the example below would be `kcmerrill.com` and `crush.kcmerrill.com`. 

```yaml

kcmerrill.com:
    url: http://kcmerrill.com
    every: 10s
    try: 3
    notify: slack 



crush.kcmerrill.com:
    cron: "*/30 * * * * *"
    shell: |
        wget -qO- https://crush.kcmerrill.com/test/something || (echo "Crush no longer accepting messages" && false)
    notify: slack
```

* *cron* is simply the same cron syntax you're used of. You can [read more about it here](https://godoc.org/github.com/robfig/cron). 

* *every* can be a golang `time.Duration` or special strings such as `@hourly`, `@daily`, `@midnight`. You can [read more about it here](https://godoc.org/github.com/robfig/cron).

* *shell* is a custom `channel`. This can be anything you'd like so long as it has a corrisponding `channel`. `spock` has a few custom `channels` built in. `slack`, `shell`, `url`. We can add more, we can also add public github repofiles as channels too! 

* *try* indicates the number of attempts `spock` should try and fail before notifying. Sometimes network connections can be finicky, scripts or whatever. By default, without `try` it will automatically alert upon it's first failure. 

* *notify* is a space separated string that describes the `channels` to send the check to. Currently, only `slack` is built in, but you can easily add `email`, `logging` or whatever suits your needs.

