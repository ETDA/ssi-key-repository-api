
<h1 align="center">
    Key Repository API
</h1>

<p align="center">
  <a href="#about">About</a> •
  <a href="#development">Development</a> 
</p>

## About

The SSI Key Repository API are the service that store encrypted private key which encryptped by Public key in HSM (Utimaco). 

## Development

### Prerequisites

- Need to connect to ETDA VPN to access to the HSM with IP `10.121.122.32` and `10.121.122.33` or change to your HSM IP Address.

#### Start Service
- Copy file `.env.sample` to `.env`
- run `docker-compose up -d`
- you can access the service via `http://localhost:8081`


