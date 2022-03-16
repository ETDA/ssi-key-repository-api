# Key Repository API

## Introduction
\
This service is using store encrypted private key which encrypted by Public Key in HSM (Utimaco)

## Step to start the service
- Copy file `.env.sample` to `.env`
- run `docker-compose up -d`
- you can access the service via `http://localhost:8081`

## Required
- need to connect to ETDA VPN to access to the HSM with IP `10.121.122.32` and `10.121.122.33`.
