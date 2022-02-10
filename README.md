# Example

This is an example go + React application which shows off the functionality of the `github.com/go-webauthn/webauthn`
library. 

## Features

This project implements the following features:

- Shows off the fact this library can work with practically any http lib (like github.com/valyala/fasthttp) thanks to 
  the hard work by the contributors to github.com/duo-labs/webauthn
- U2F Registration
- Attestation:
  - Standard
  - Resident Key / Discoverable
- Assertion:
  - Standard
  - AppID (U2F Registered Keys)
  - Discoverable (Passwordless)
- Simulates Username/Password login via a username only login.
- Shows details about browser support.
- Shows debug messages relevant to what occurred; for example if the user chooses cancel.
- Full TypeScript implementation of the Webauthn portions.

## Requirements

The following requirements are needed to build the project:

- npm
- go (1.17)

## Steps

The following steps allow you to run the project:

- Copy `config.example.yml` to `config.yml`.
- Edit `config.yml` with your desired settings.
- Run `go generate` to generate the react frontend.
- Run `go build ./cmd/server`.
- Run the `./build` binary.
- Configure a HTTPS proxy to answer requests on your desired domain.
- Visit the domain.

