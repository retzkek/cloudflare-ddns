cloudflare-ddns
===============

Dynamic DNS client using CloudFlare's client 
[API](http://www.cloudflare.com/docs/client-api.html). 

Copyright 2013 Kevin Retzke <kmr@kmr.me>

Shared under the MIT license, see `LICENSE` for details.

Installing
----------
You will need the [Go](http://golang.org) development environment installed
to build/run this program. Only packages in the standard library are used.

After that is set up:

    go get github.com/retzkek/cloudflare-ddns

Set `EMAIL`, `ZONE`, and `DOMAIN` appropriately (for most cases `DOMAIN` == `ZONE`, 
unless you have a dynamic subdomain.

`TKN` is your CloudFlare API key from https://www.cloudflare.com/my-account.

`EGRESS` is the publicly-facing network interface. If you want to run this on a 
machine within your network, the getAddr() function will need to be changed
to optain the IP address via some external means (e.g. web scraping). Let me know
if you need this.

`FREQUENCY` is how often you want to check for IP address changes. Note that it will
only hit the CloudFlare API when the IP address changes, so it doesn't hurt to make it
fairly often.

After setting your options, change to the package directory 
(`$GOPATH/pkg/src/github.com/retzkek/cloudflare-ddns`) and run:

    go install
    
This will build the executable and install it in `$GOPATH/bin`.

Usage
-----
If `$GOPATH/bin` is in your path then you just have to run:

    cloudflare-ddns

This will run in the foreground with logging to stderr.  You could just
run in a tmux or screen session, or background it:

    nohup cloudflare-ddns 2>log.file &

Bugs/Requests
-------------
Feel free to email me <kmr@kmr.me> or submit an issue on GitHub.
