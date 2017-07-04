# Logvault

A simple and centrallized logging service

## Is it working?

Yes, the service is working with very minimal features as we want to use logvault as fast as possible.

For now, the log file is being kept in logvault service server in a directory and in a form of `*.log` file

## Who is logvault for?

- Anyone that can bear with this experimental project
- Anyone who is still saving log in log files and save the logs in local directory.
- Anyone who want a simple centralized log service where you can find your logs easily.

## How logvault works?

Logvault agent need to be placed and running in the server that you want to ingest its log files. The agent can ingest a spesific log file or entire directory.

Agent is a very small golang binary that can be placed in the server.

All Ingested logs will sent to a centralized logvault service, where all the logs being kept.

## TODO
1. Alert rules from logs.
2. Data structure for logs.
3. Search from structured logs.