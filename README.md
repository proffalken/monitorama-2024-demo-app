1. Make sure you've got alloy running on localhost, accepting OTEL data via HTTP/Protobuf, and pointing at a Grafana Cloud stack
2. Clone https://github.com/proffalken/monitorama-2024-demo-app
3. Create a new virtual env and run `pip install -r requirements.txt`
4. Execute `./bin/launch.sh`
5. Run the following: `cd basic_otel && ./manage.py createsuperuser`, then fill out the prompts
6. Log in to [http://localhost:8898/admin](http://localhost:8898/admin) and add a couple of spaces based on the [SpaceAPI Directory](https://github.com/SpaceApi/directory/blob/master/directory.json)
7. Run `./bin/launch_weather.sh`
8. Run `./bin/launch_client.sh`

You should start to see the traces and metrics flow through the platform and into Grafana Cloud App O11y. 

## Containers

Two Docker images are provided:

- `basic_otel/Dockerfile`
- `space_weather/Dockerfile`

Both images run database migrations on startup and launch each Django service via `opentelemetry-instrument`.

Local image build examples:

```bash
docker build -f basic_otel/Dockerfile -t basic-otel:dev .
docker build -f space_weather/Dockerfile -t space-weather:dev .
```

## GitHub Image Workflow

The workflow at `.github/workflows/build-python-images.yml` builds both Python service images and pushes them to GHCR on `main` and tags.

Published image names:

- `ghcr.io/<org-or-user>/monitorama-2024-demo-app-basic-otel`
- `ghcr.io/<org-or-user>/monitorama-2024-demo-app-space-weather`

## Helm Deployment

Helm chart path:

- `deploy/charts/monitorama-demo`

Quick start:

```bash
helm upgrade --install monitorama-demo deploy/charts/monitorama-demo \
  --set basicOtel.image.repository=ghcr.io/<org-or-user>/monitorama-2024-demo-app-basic-otel \
  --set spaceWeather.image.repository=ghcr.io/<org-or-user>/monitorama-2024-demo-app-space-weather \
  --set basicOtel.image.tag=latest \
  --set spaceWeather.image.tag=latest \
  --set spaceWeather.weatherApiKey=<your-weatherapi-key>
```

By default, the chart configures OTLP export to `http://otel-collector:4318` with `http/protobuf`. Override `global.otel.*` values as needed.

### CloudNativePG

The Helm chart includes native CloudNativePG support using the CNPG operator CRD (`postgresql.cnpg.io/v1`, `Cluster`).

Enable it like this:

```bash
helm upgrade --install monitorama-demo deploy/charts/monitorama-demo \
  --set cnpg.enabled=true \
  --set cnpg.database.password='<strong-password>' \
  --set basicOtel.image.repository=ghcr.io/<org-or-user>/monitorama-2024-demo-app-basic-otel \
  --set spaceWeather.image.repository=ghcr.io/<org-or-user>/monitorama-2024-demo-app-space-weather
```

When enabled, the chart:

- Creates a CNPG `Cluster`
- Uses CNPG bootstrap (`initdb`) to create the app database and owner
- Wires `basic_otel` to the CNPG read-write service and app credentials

If you already manage DB credentials externally, set `cnpg.database.existingSecret` and the chart will use that secret instead of creating one.
