# Wireguard Manager And API

A manager and API to add, remove clients as well as other features such as an auto reapplier which deletes and adds back a client after inactivity to increase their privacy by removing their IP address from memory.

This GoLang application runs an API which can be made **https** ready using a LetsEncrypt certificate. The program creates directories in the directory ``/opt/wgManagerAPI`` (This needs to be created manually before hand). In the ``/opt/wgManagerAPI`` directory we have a few more sub-directories such as ``/logs`` which contain logs of the application and ``/wg`` which contains our SQLite database.

The SQLite database contains tables which store information such as generated and available IPs, client configuration (public key and preshared key) as well as the Wireguard server own private key, public key, IP Addresses and ListenPort.

## Content
- [Wireguard Manager And API](#wireguard-manager-and-api)
  - [Content](#content)
  - [Config.json File](#configjson-file)
    - [Instance settings](#instance-settings)
    - [API server settings](#api-server-settings)
  - [Deployment](#deployment)
    - [Docker](#docker)
      - [Docker Compose](#docker-compose)
    - [Building from source](#building-from-source)
      - [Code](#code)
      - [Dockerfile](#dockerfile)
  - [Communicating with the API](#communicating-with-the-api)
    - [Adding keys](#adding-keys)
    - [Removing keys](#removing-keys)
    - [Enabling keys](#enabling-keys)
    - [Disabling keys](#disabling-keys)
    - [Getting all keys](#getting-all-keys)
    - [Editing subscriptions](#editing-subscriptions)
  - [Debugging](#debugging)
    - [Logs](#logs)
    - [FAQ](#faq)
## Config.json File
A config.json file needs to be placed in the directory `/opt/wgManagerAPI/config.json`. A template can be found in the `src/config/template.json`.

### Instance settings
| Variable | Purpose | Type |
| ------------ | ------------ | ------------ |
| INSTANCE.IP.GLOBAL.ADDRESS.IPV4  | The public IPv4 **addresses** of your server.  Must be set.| string array |
| INSTANCE.IP.GLOBAL.ADDRESS.IPV6  | The public IPv6 **addresses** of your server.  Must be set.| string array |
|  INSTANCE.IP.GLOBAL.DNS | The DNS address that you want wireguard clients to connect to. Can also be a local address if you are running a Pihole instance or local DNS.  | string |
|  INSTANCE.IP.GLOBAL.ALLOWED |  By default it allows all IPv4 and IPv6 addresses through. Change to allow split tunneling. Default of ``0.0.0.0/0, ::0``| string |
|  INSTANCE.IP.LOCAL.IPV4.ADDRESS |  The local IPv4 address which will be assigned to the Wireguard instance. **IMPORTANT:** the application creates a subnet of /16, please make sure you have space for this. By default it is set to ``10.6.0.1`` (p.s. this was tested with a Pihole instance running locally on the same address).  |  string |
|  INSTANCE.IP.LOCAL.IPV4.SUBNET |  Subnet of the local IPv4 address.  |  string |
|  INSTANCE.IP.LOCAL.IPV6.ADDRESS |  The local IPv6 address which will be assigned to the Wireguard instance. **IMPORTANT:** At the current stage the docker container is not able to access IPv6, only IPv4. If you would like to disable/not use IPv6, set this to ``-``  | string |
|  INSTANCE.IP.LOCAL.IPV6.SUBNET |  Subnet of the local IPv6 address.  | string |
|  INSTANCE.IP.LOCAL.IPV6.ENABLED |  Enabling of IPv6 (does not work with docker) |  boolean |
|  INSTANCE.PORT |  Wireguard server port. Default value of 51820.  | integer |

### API server settings
| Variable | Purpose | Type |
| ------------ | ------------ | -----------|
|  SERVER.MAX_IP | The number of IPs that will be generated in the SQLite database as well as the maximum number of clients that the server can host  | string |
|  SERVER.SECURITY |  Enables HTTPS on the server. A FULLCHAIN_CERT and PK_CERT must be specified. Set to ``disabled`` to use a HTTP connection and anything else to enable HTTPS. By default is set to true.| boolean |
|  SERVER.CERT.FULLCHAIN | The path to your LetsEncrypt fullchain.pem certificate. For example: ``/etc/letsencrypt/live/domain.com/fullchain.pem`` | string |
|  SERVER.CERT.PK | The path to your LetsEncrypt privkey.pem certificate. For example: ``/etc/letsencrypt/live/domain.com/privkey.pem``  | string|
| SERVER.AUTH  | The Authorisation key that is needed in an API request ``Authentication`` header. Setting this to a ``-`` will disable API authentication. Default value of "ABCDEFG" | string |
| SERVER.PORT  | The port that is used to run the API server (this is not the Wireguard server port). Default of port of 8443 | string |
| SERVER.AUTOCHECK | Enable the autochecker (automatically deletes and re-adds client keys after inactivity to increase privacy of user) by setting this to ``enabled``. Disable by setting to ``-``.| boolean |
|  SERVER.INTERFACE | The interface of your network card. Usually eth0.| string |

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

The docker-compose file is the easiest way to get software up and running. Do not forget to add your ``config.json`` file to ``/opt/wgManagerAPI/config.json``

### Building from source
#### Code
Building from source allows you to create an executable file which can be created into a Systemd service or equivalent. It also allows you to build for a different architecture such as ARM64. Running the executable must be run with sudo (recommended) or root (not recommended). 

Do not forget to add your `config.json` file to `/opt/wgManagerAPI/config.json`
1. Install Go 1.14+ on to your machine
2. git clone this repository
3. ``cd wireguard-manager-and-api`` to open the repo
4. ``go get`` to get packages
5. ``go build -o wgManagerAPI main.go`` to build an output a executable file
6. ``sudo ./wgManagerAPI`` to run the application.

#### Dockerfile
Building a docker image from scratch enables you to create an image specific to a specific architecture such as ARM64 as prebuilt images in our docker image repository is made for AMD64 architecture.

Do not forget to add your ``config.json`` file to ``/opt/wgManagerAPI/config.json``
1. Install docker
2. Clone the git repository and open the directory
3. ``sudo docker build -t wireguard-manager-and-api:YOURTAGHERE .`` to build the docker image locally

## Communicating with the API

With almost any API error the server will give back a ``400 Bad Request`` status code. Please read the JSON response file "response" to get the error information.

### Adding keys

URL: `POST` request to `http(s)://domain.com:PORT/manager/key`  
Header: `Content-Type: application/json`  
Header (If authentication is enabled): `authorization:(AUTH key from config.json)`  
Body: 
```json
{
  "publicKey": "(Wireguard client public key)",
  "presharedKey": "(Wireguard client preshared key)",
  "bwLimit": (integer, megabytes, 0 for infinite),
  "subExpiry": "(date in this FORMAT!! uses UTC for timing) 2021-Oct-28 12:39:05 PM",
  "ipIndex": "(the integer index of the ip address you want to use. Index is the IP in the config.json on the API server.)
}
```
Response:  
Status Code `202`
```json
{
  "allowedIPs": "(Allowed IPs set in config.json. Default of 0.0.0.0/0, ::0)",
  "dns": "(DNS set in config)",
  "ipAddress": "(public IP of server)",
  "ipv4Address": "(internal IPv4 Address assigned to client) 10.6.0.10/32", 
  "ipv6Address": "(internal IPv6 Address assigned to client) fe22:22:22::10/128",
  "keyID": "(KeyID in database) 1",
  "listenPort": "(wireguard default listenport) 51820",
  "publicKey": "(Public key of wireguard server) ghyewr34A0wzT1b7ZdJgPjWwS3F/9PgRzlNWcX/QlA0=",
  "response": "Added key successfully"
}
```
Parsing response into Wireguard client config:
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
URL: `DELETE` request to `http(s)://domain.com:PORT/manager/key`  
Header: `Content-Type: application/json`  
Header (If authentication is enabled): `authorization:(AUTH key from config.json)`  
Body:
```json
{
  "keyID": "(Database keyID)"
}
```
Response:  
Status Code `202`
```json
{
  "response": "Key deleted"
}
```

### Enabling keys
URL: `POST` request to `http(s)://domain.com:PORT/manager/key/enable`  
Header: `Content-Type: application/json`  
Header (If authentication is enabled): `authorization:(AUTH key from config.json)`  
Body:
```json
{
  "keyID": "(Database keyID)"
}
```
Response:  
Status Code `202`
```json
{
  "response": "Enabled key successfully"
}
```

### Disabling keys
URL: `POST` request to `http(s)://domain.com:PORT/manager/key/disable`  
Header: `Content-Type: application/json`  
Header (If authentication is enabled): `authorization:(AUTH key from config.json)`  
Body:  
```json
{
  "keyID": "(Database keyID)"
}
```
Response:  
Status Code `202`
```json
{
  "response": "Disabled key successfully"
}
```
### Getting all keys
This requests all keys from the database. It is basically getting all information from the Keys table.

URL: `GET` request to `http(s)://domain.com:PORT/manager/key`  
Header (If authentication is enabled): `authorization:(AUTH key from config.json)`

Response:  
Status Code `202`
```json
{
  "response": "All key successfully parsed",
  "Keys": [
    {
      "KeyID": 1,
      "PublicKey": "Si0sJ6YoYOvyc96Ln6ZkBYRXOz08BNJfM3iFjoAiyU=",
      "PresharedKey": "AEDWepAAXPzWzatzShkKKJ4UaWscDMetlDq4d0OXfp8=",
      "IPv4Address": "10.6.0.10",
      "Enabled": "true"
    },
    {
      "KeyID": 2,
      "PublicKey": "0ba1shYHi9swnUvawqXuVDSM//S9OW2KxdvoaF69NHg=",
      "PresharedKey": "dKVFJlFK387Ht174ol0eeW/gjVfZZGLIfSN4egInLgY=",
      "IPv4Address": "10.6.0.100",
      "Enabled": "true"
    },
    ...
  ]
}
```

### Editing subscriptions 
This allows editing of subscriptions such as bandwidth, resetting bandwidth usage and changing the subscription expiry date.  
URL: `POST` request to `http(s)://domain.com:PORT/manager/subscription/edit`  
Header: `Content-Type: application/json`  
Header (If authentication is enabled): `authorization:(AUTH key from config.json)`  
Body:  
```json
{
	"keyID": "(database keyID)",
	"bwLimit": (integer, bandwidth limit. Set to -1 to keep current limit.),
	"subExpiry": "(date, must be this format!!) 2032-Oct-21 12:49:05 PM",
	"bwReset": true (a boolean, set to true to reset current bandwidth usage.)
}
```
Response:  
Status Code `202`
```json
{
  "response": "Updated successfully"
}
```

## Debugging
### Logs
If the Wireguard Manager and API application fails to start you should always look at your logs and the errors to see the problems look at ``/opt/wgManagerAPI/logs/`` folder and open the latest log using ``nano`` or any other text editor.

### FAQ
**Q:** The prebuilt source file or docker image is not working properly.  
**A:** Build from dockerfile or code from source. The prebuilt docker images are not for ARM architecture.

**Q:** IPv6 is not working with the docker image.  
**A:** We have not been able to setup IPv6 on the docker-compose file successfully. If you find a solution please tell us.

**Q:** I have built from source code but unable to successfully route clients through the VPN  
**A:** You may need the iptables rule: ``sudo iptables -t nat -A POSTROUTING -o (INTERFACE I.E eth0 or enp0s3) -j MASQUERADE``. This will be required on boot everytime. We will try and implement this into the program in the future.

**Q:** Do I need to use a  wireguard (wg0.conf) file?  
**A:** No, and please do not try, it may mess up some of the functionality we provide such as automatic deleting and re-adding keys.



