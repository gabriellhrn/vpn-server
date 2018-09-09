# vpn-server

Command line tool to spin-up a VPN server in Digital Ocean.


## Requirements

  * go1.11+


## How to use

Export `DO_PAT` - the Digital Ocean Personal Authentication Token -, and `DO_SSHKEY` - fingerprint of a SSH key configured in Digital Ocean:

```
$ set -x DO_PAT <token>
$ set -x DO_SSHKEY <fingerprint>
```

After configuring the environment variables, *start* the server:

```
$ vpn-server start
```

To *stop* the server, run:

```
$ vpn-server stop
```

## Further improvements

  * Automatically enable Digital Ocean Cloud Firewall
  * Automatically connect to the VPN after the server is created
  * Configure peer in droplet user-data
