# TUNNEL-ALIVE

## Description
A simple cli-app which allows SSH tunnelling(Port forwarding)

## Motivation
When using MobaXterm for SSH tunnelling it only allows 3 saved tunnels for the personal edition.
But during my development I need to open more than 3 sessions at a time
## How it works
First it reads all session info in a config file, then for each session
it open a SSH tunnelling according to the configuration.
For preventing session from being killed after long idle time I periodically send null message to the remote server.
## How to run
`go run main.go`