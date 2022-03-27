# dns-study

This repo is a study on creating a custom DNS server. I wanted to go through the
implementation and specification document [RFC1035] and create it using GoLang
to get a better understanding on how DNS works. I'm relatively new to Golang so 
this is also test to see what I can do and learn with this language.

## Usage

```bash
# Get the soruce code
git clone git@github.com:lcox74/dns-study.git
cd dns-study

# Build
go build

# Run the program
./dns-study.exe
```

As this is a DNS server, you will have to change in either your browser settings
or your IPv4 Network Adapter Options to use this as your DNS. On Windows you'll
have to change the DNS server address to:

[Windows DNS Settings](docs/res/windowsdns.PNG)

## What currently works?

The application currently works as a DNS Proxy which listens on the DNS port
`:53` and then redirects any requests to the CloudFlare (`1.1.1.1`) DNS server
and forwards the result back to the client. This is so I can test if I am able
to extract the bytes into the proper structures.

[RFC1035]: https://datatracker.ietf.org/doc/html/rfc1035