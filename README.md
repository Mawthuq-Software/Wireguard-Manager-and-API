# Wireguard Manager And API

A manager and API to add, remove clients as well as other features such as an auto reapplier which deletes and adds back a client after inactivity to increase their privacy by removing their IP address from memory.

This GoLang application runs an API which can be made **https** ready using a LetsEncrypt certificate. The program creates directories in the directory ``/opt/wgManagerAPI`` (This needs to be created manually before hand). In the ``/opt/wgManagerAPI`` directory we have a few more sub-directories such as ``/logs`` which contain logs of the application and ``/wg`` which contains our SQLite database.

The SQLite database contains tables which store information such as generated and available IPs, client configuration (public key and preshared key) as well as the Wireguard server own private key, public key, IP Addresses and ListenPort.
# How to use
## .env File
A .env file needs to be placed in the directory `/opt/wgManagerAPI/.env` containing the following:

```bash
MAX_IP=350
SERVER_SECURITY=enabled
FULLCHAIN_CERT=
PK_CERT=
AUTH=ABCDEFG
IP_ADDRESS=
DNS=1.1.1.1
ALLOWED_IP=0.0.0.0/0, ::/0

WG_IPV4=10.6.0.1
WG_IPV6=fe22:22:22::1
PORT=8443
AUTOCHECK=enabled
```
| Variable  | Purpose  |
| ------------ | ------------ |
|  MAX_IP | The number of IPs that will be generated in the SQLite database as well as the maximum number of clients that the server can host  |
|  SERVER_SECURITY |  Enables HTTPS on the server. A FULLCHAIN_CERT and PK_CERT must be specified. Set to ``disabled`` to use a HTTP connection and anything else to enable HTTPS. |
|  FULLCHAIN_CERT | The path to your LetsEncrypt fullchain.pem certificate. For example: ``/etc/letsencrypt/live/domain.com/fullchain.pem`` |
|  PK_CERT | The path to your LetsEncrypt privkey.pem certificate. For example: ``/etc/letsencrypt/live/domain.com/privkey.pem``  |
| AUTH  | The Authorisation key that is needed in an API request ``Authentication`` header. Setting this to a ``-`` will disable API authentication  |
| IP_ADDRESS  | The public IP address of your server.  |
|  DNS | The DNS address that you want wireguard clients to connect to. Can also be a local address if you are running a Pihole instance or local DNS.  |
|  ALLOWED_IP |  By default it allows all IPv4 and IPv6 addresses through. Change to allow split tunneling. |
| WG_IPV4  | The local IPv4 address which will be assigned to the Wireguard instance. **IMPORTANT:** the application creates a subnet of /16, please make sure you have space for this. By default it is set to ``10.6.0.1`` (p.s. this was tested with a Pihole instance running locally on the same address).|
| WG_IPV6  | The local IPv6 address which will be assigned to the Wireguard instance. **IMPORTANT 1.1:** the application creates a subnet of /64, please make sure you have space for this. By default it is set to ``fe22:22:22::1``  **IMPORTANT 1.2:** At the current stage the docker container is not able to access IPv6, only IPv4. If you would like to disable/not use IPv6, set this to ``-``.|
| PORT  | The port that is used to run the API server (this is not the Wireguard server port). |
|  AUTOCHECK  | Enable the autochecker (automatically deletes and re-adds client keys after inactivity to increase privacy of user) by setting this to ``enabled``. Disable by setting to ``-``.|

## Deployment
### Docker

A docker container is automatically built on a new release. For this repository, the [container registry](https://gitlab.com/mawthuq-software/wireguard-manager-and-api/container_registry/2171069 "container registry") has tags relevant to the docker image. The **main** tag refers to a stable release and **latest** refers to a newly built image. This may be unreleased or buggy software so use the latest tag with caution.

Our docker image is built with Debian buster and CoreDNS is used to allow the internal docker container DNS to communicate with the host DNS. 

**IMPORTANT:** Currently with the Docker setup IPv6 addresses cannot passthrough, only IPv4 addresses.
#### Docker Compose
```yaml
version: "3"

services:
    wireguard-manager-and-api:
      image:  registry.gitlab.com/mawthuq-software/wireguard-manager-and-api:main
      volumes:
      - /etc/letsencrypt:/etc/letsencrypt
      - /opt/wgManagerAPI:/opt/wgManagerAPI
      - /lib/modules:/lib/modules
      ports:
      - "8443:8443"
      - "51820:51820/udp"
      cap_add:
        - NET_ADMIN
        - SYS_MODULE
      sysctls:
        - net.ipv4.conf.all.src_valid_mark=1
        - net.ipv6.conf.all.disable_ipv6=0
```

The docker-compose file is the easiest way to get software up and running. Do not forget to add your ``.env`` file to ``/opt/wgManagerAPI/.env``

### Building from source
Building from source allows you to create an executable file which can be created into a Systemd service or equivalent. Running the executable must be run with sudo (recommended) or root (not recommended). 

Do not forget to add your ``.env`` file to ``/opt/wgManagerAPI/.env``
1. Install Go 1.14+ on to your machine
2. git clone this repository
3. ``cd wireguard-manager-and-api`` to open the repo
4. ``go get`` to get packages
5. ``go build -o wgManagerAPI main.go`` to build an output a executable file
6. ``sudo ./wgManagerAPI`` to run the application.

## Communicating with the API
### Adding keys

URL: `POST` request to `http(s)://domain.com:PORT/manager/keys`

Header: `Content-Type: application/json`

Header (If authentication is enabled): `authorization:(AUTH key from .env)`

Body: 
```json
{
  "publicKey": "(Wireguard client public key)",
  "presharedKey": "(Wireguard client preshared key)"
}
```
Response: 
```json
{
  "allowedIPs": "0.0.0.0/0, ::/0",
  "dns": "10.6.0.1",
  "ipAddress": "(public IP of server)",
  "ipv4Address": "(internal IPv4 Address assigned to client) 10.6.0.10/32", 
  "ipv6Address": "(internal IPv6 Address assigned to client) fe22:22:22::10/128",
  "keyID": "(KeyID in database) 1",
  "listenPort": "(wireguard default listenport) 51820",
  "publicKey": "(Public key of wireguard server) ghyewr34A0wzT1b7ZdJgPjWwS3F/9PgRzlNWcX/QlA0=",
  "response": "Added key successfully"
}
```
Parsing response into client config:
```
[Interface]
PrivateKey = (Wireguard client private key)
Address = (internal IPv4 Address assigned to client), (internal IPv6 Address assigned to client)
DNS = 10.6.0.1

[Peer]
PublicKey = (Public key of wireguard server)
PresharedKey = (Wireguard client preshared key)
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = (public IP of server):51820
```

### Removing keys
URL: `DELETE` request to `http(s)://domain.com:PORT/manager/keys`

Header: `Content-Type: application/json`

Header (If authentication is enabled): `authorization:(AUTH key from .env)`

Body:
```json
{
  "keyID": "(Database keyID)"
}
```
Response:
```json
{
  "response": "Key deleted"
}
```


## Debugging
### Logs
If the Wireguard Manager and API application fails to start you should always look at your logs and the errors to see the problems look at ``/opt/wgManagerAPI/logs/`` folder and open the latest log using ``nano`` or any other text editor.

### FAQ
Haha nothing here




