# GoDDNS - Go Dynamic DNS server

[![go.dev reference](https://img.shields.io/badge/go-reference-blue)](https://pkg.go.dev/github.com/svenhertle/goddns)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=svenhertle_goddns&metric=alert_status)](https://sonarcloud.io/dashboard?id=svenhertle_goddns)

GoDDNS is a dynamic DNS updater.
It receives DNS updates from clients, authenticated via a secret key, and writes the changes to an existing DNS server.

## Why another one?

* Support for IPv4 and IPv6
* Extensible architecture: new backends (= DNS servers) can be added easily
* Why not?

## Backends

Currently, the following backends are supported:

* PowerDNS with PostgreSQL
* Dummy (just for testing/debugging)

# Use it

## Install

It is expected here that the DNS server runs already, check the [Wiki](https://github.com/svenhertle/goddns/wiki) for instructions.

GoDDNS is just one binary and a config file. You can run it with or without reverse proxy in front of it.
It supports TLS natively, so there is no need for a reverse proxy.
This short tutorial explains how to install GoDDNS to `/srv/goddns`, but feel free to change this.

Do not run GoDDNS as root, add an dedicated user instead:

    adduser --disabled-login --disabled-password --home /srv/goddns goddns

You can download the latest [release from Github](https://github.com/svenhertle/goddns/releases) or build it with `make`.
Place it in `/srv/goddns/goddns-x64-linux` and copy the [configuration file](https://github.com/svenhertle/goddns/blob/master/goddns.yaml) to `/srv/goddns/goddns.yaml`.
The location of the configuration file can also be `/etc/goddns/goddns.yaml` or specified via command line.
Change the configuration file according to your needs, the available options are documented directly in the file.

Configure the reverse proxy (see configuration for [Apache](https://github.com/svenhertle/goddns/blob/master/install/apache.conf) for example) or TLS directly in GoDDNS.

The [systemd service file](https://github.com/svenhertle/goddns/blob/master/install/goddns.service) should be copied to `/etc/systemd/system/goddns.service`.
If necessary, adapt it to your needs and then enable it:

    systemctl daemon-reload
    systemctl enable --now goddns

After the GoDDNS configuration was changed, the service needs to be restarted:

    systemctl restart goddns

## API for clients

Option 1: Specify IP address as parameter:

    https://dyn.example.org/v1/?key=foo&ip=1.2.3.4

Option 2: Do not specify the IP address explicitly, use the client IP address visible for GoDDNS:

    https://dyn.example.org/v1/?key=foo

Check the [Wiki](https://github.com/svenhertle/goddns/wiki) for hints how to use specific clients.

## HTTP response codes

The possible HTTP response codes are:

| Code     | Text                |
|----------|---------------------|
| 200      | ok                  |
| 400      | key missing         |
| 400      | key unkown          |
| 400      | ip invalid          |
| 500      | internal error      |

# Development

## Adding a new backend

One design goal of GoDDNS is that new backends can be easily added.
All backends need to implement the interface `Backend` (see [backend.go](`https://github.com/svenhertle/goddns/blob/master/backends/backend.go`)).

* Copy `backends/dummy.go` to new `backends/newbackend.go`
* Adapt file: rename `DummyBackend` to new name
* Add new backend in `backends/backend.go` in `GetBackend` method
* Implement `Configure` and `Update` methods.

`Configure` receives a Viper configuration, so custom configuration options can be implemented directly.
The [PowerDNS backend](https://github.com/svenhertle/goddns/blob/master/backends/powerdns.go) in `powerdns.go` can be used as an example.

## Release new version

Update `goddnsVersion` in `main.go` and update Git repo:

    git commit -am "bump version"
    git tag v1.2.3
    git push
    git push --tags

Build new binary with `make` and upload a new release to Github.

# License

GNU General Public License v3.0 or later

See [LICENSE](https://github.com/svenhertle/goddns/blob/master/LICENSE) to see the full text.