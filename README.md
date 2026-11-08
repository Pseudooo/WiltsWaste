 
### Overview

Small web application built with svelte & go that allows users to subscribe to a live calendar stream for their bin days, all information
used to drive the calendars is collected from the official [wiltshire.gov.uk](https://www.wiltshire.gov.uk/rubbish-and-recycling) site.

This site is publicly available at [wiltswaste.uk](https://wiltswaste.uk/)

### Running locally

#### Docker Compose

If you have docker installed on your system then the docker compose file can be used to run the application locally, simply running:
```
docker compose up --build
```
In the repository root directoy will make the website available at `http://localhost:8080`

#### API

The API is built with with go & gin, so the go CLI must be installed to compile it. Alternatively the application is containerised so docker
or podman can be used instead

```
cd ./api/
go build .
```

Which will create `WiltsWaste.exe` ready to run


##### Monitoring

This application is setup for structured json logging and open telemetry for metrics. Logs are exported
to `STDOUT`.

To enable metrics the following environment variables must be set:
- `OTEL_EXPORTER_OTLP_ENDPOINT`
- `SERVICE__VERSION`

Otherwise the application will not start.

#### Web

The web project is a SPA built with vite, svelte & shadcn, with node & npm installed just use the normal `npm run dev` or make use of the dockerfile.

The API base url needs to be declared in the `VITE_API_BASE_URL` env var
