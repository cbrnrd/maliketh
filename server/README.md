# Maliketh - Server

## Features

* Multi-user (operators)
* Easily configurable (via YAML files)
* Easily deployable (via Docker)
* Per-operator implant builder


## Setup

To start the server, 90% of your work can be done by running the following command in the project root:

```bash
docker-compose -f server/docker-compose.yml --env-file server/.env.example up
```

Note: You will need to create a `.env` file in the `server/` directory. See `.env.example` for an example.

The only thing left to do is bootstrap the database and create the admin user. To do this, run the following command in the `server` directory:

```bash
./bootstrap_db.sh
```

The output of this script will be a JSON configuration for the admin user. You can use this with the maliketh [client](../client/) to connect to the server.


## Ideal server setup

An ideal setup would involve 2 servers. 1 running nginx which the implants connect back to, and 1 running the actual server. This would allow you to use a domain name for the implants to connect to, and also allow you to use SSL. The nginx server would be configured to proxy all traffic to the server. The nginx server would also be configured to use SSL. The server would be configured to only accept connections from the nginx server. This would allow you to use SSL, but not have to worry about the overhead of SSL on the server.

On the server side, Wireguard should be installed and unique WireGuard keys should be given to each operator along with their operator configuration. The server should be configured to only accept connections from the Wireguard interface. Wireguard keys should be generated for each operator. The server should be configured to only accept connections from the Wireguard interface.

<p align="center">
  <img src="./data/Maliketh%20Network%20Diagram.png" alt="Ideal setup" width="500"/>
</p>
