# dumbdumb

Make your dumbphones kinda smart. Service that accepts text message queries via
emails to built-in SMTP server, delegates to the appropriate APIs to gather
information and responds succinctly via text messages (email) to the sender.


## End-user Interface

### "Weather" handler

Provides weather forcast for today, tonight, and tomorrow via three text
messages.

```
[Ww]eather <town descriptor>
```

Example query:

```
weather boston, ma
```

### "Find" handler

Finds the closest matching place given an open-ended query and returns
information about it.

```
[Ff]ind <query>
```

Example query:

```
Find thrift store west roxbury, ma
```

Example response:

```
Savers
Phone: 1 617-323-8231
Addr: 1230 Veterans of Foreign Wars Pkwy, West Roxbury, MA
Open now: yes
```

### "Translate" handler

Translates a message from a source to target language.

```
[Tt]ranslate <src-lang> -> <target-lang> <message>
```

Example query:

```
Translate es -> en hola
```

Example response:

```
[es] -> [en]: Hello
```


## Prerequisites

* Server with port 25 exposed to public, with DNS A-record or MX-record
  pointing to it
* Mandrill (MailChimp) account and API key handy, and hostname for above server
  verified as a sender in the Mandrill interface (Hint: start up the `dumbdumb`
  service and watch the output to receive your Mandrill verification email.)
* Weather Underground account and API key
* Google Cloud Platform API key with Places and Translate services enabled


## Installation (Ubuntu server)

These are steps to install, tested on an ubuntu 14.04 server.

1) Install required packages and Golang:

```bash
sudo apt-get update
sudo apt-get install -y git mercurial curl
sudo curl https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz \
  | tar -C /usr/local -xz
```

2) Set GOPATH and ensure it is on PATH:

```bash
echo "export GOPATH=\$HOME/go" >> $HOME/.bash_profile
echo "export PATH=\$PATH:/usr/local/go/bin:\$HOME/go/bin" >> $HOME/.bash_profile
```

3) Clone this repo into `$GOPATH/src/github.com/davemt/dumbdumb`

4) Fetch dependencies and install the service entry point script:

```bash
cd $GOPATH/src/github.com/davemt/dumbdumb
go get ./...
cd dumbdumb && go install
```


## Running the service

Ensure all required environment variables are set (see `settings.sh.example`).

Start up service with `dumbdumb` (requires `$GOPATH/bin` to be on `$PATH`).

It is not a daemon so you'll need to do any detaching or output redireciton on
your own.
