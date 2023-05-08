<!-- Centered logo -->
<p align="center">
  <img src="./data/maliketh_logo.png" alt="Maliketh logo" width="900" height="300"/>
</p>

# Maliketh
Maliketh is a multi-user, customizable C2 framework. The goal of Maliketh is to provide a flexible, easy to use C2 framework that can be customized to fit the needs of the operator. The poster used in the initial presentation is located [here](./data/Maliketh%20C2%20Poster.png).

## Server features
* Multi-user (operators)
* Easily configurable (via YAML files)
* Easily deployable (via Docker)

## Implant features
The implant is written in C++ and targeted for Windows. The main feature of the implant is its ability to change its behavior based on the configuration file it receives from the server. This allows the operator to customize the implant to fit their needs. The implant also has the following features (see [here](./design/opcodes.md) for more info):
* File upload/download
* Command execution
* Shellcode injection
* Update configuration
* Send system information
* Self-destruct
* Sleep
* Basic Anti-debugging
* *Very* Basic Anti-VM
* Sleep skipping detection


## Server deployment
To start the server, 90% of your work can be done by running the following command in the `server/` directory:

```bash
docker-compose --env-file .env up -d
```

Note: You will need to create a `.env` file in the `server/` directory. See `.env.example` for an example.

The only thing left to do is bootstrap the database and create the admin user. To do this, run the following command:

```bash
./bootstrap_db.sh
```

The output of this script will be a JSON configuration for the admin user. You can use this with the maliketh [client](./client/) to connect to the server.

## Ideal server setup
An ideal setup would involve 2 servers. 1 running nginx which the implants connect back to, and 1 running the actual server. This would allow you to use a domain name for the implants to connect to, and also allow you to use SSL. The nginx server would be configured to proxy all traffic to the server. The nginx server would also be configured to use SSL. The server would be configured to only accept connections from the nginx server. This would allow you to use SSL, but not have to worry about the overhead of SSL on the server.

On the server side, Wireguard should be installed and configured. The server should be configured to only accept connections from the Wireguard interface. Wireguard keys should be generated for each operator. The server should be configured to only accept connections from the Wireguard interface.

<p align="center">
  <img src="./data/Maliketh%20Network%20Diagram.png" alt="Ideal setup" width="500"/>
</p>

## Future work
- [ ] Implement Golang client
- [ ] Per-operator builder in-server
- [ ] Stealer/basic looter
- [ ] Keylogger
- [ ] Route RabbitMQ traffic through Admin listener instead of directly connecting
- [ ] Improved anti-vm (check BIOS information)
- [x] More stable file uploads/downloads ([91a40f2](https://github.com/cbrnrd/maliketh/commit/91a40f2ba1cded5a025004a6143578fa84baec66))
- [ ] Alternate C2 channels (WireGuard, DNS, Discord, Slack, etc.)
- [x] Basic OS functions built in ([91a40f2](https://github.com/cbrnrd/maliketh/commit/91a40f2ba1cded5a025004a6143578fa84baec66))
- [x] Situational Awareness ([91a40f2](https://github.com/cbrnrd/maliketh/commit/91a40f2ba1cded5a025004a6143578fa84baec66))
