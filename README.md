1. Make sure you've got alloy running on localhost, accepting OTEL data via HTTP/Protobuf, and pointing at a Grafana Cloud stack
2. Clone https://github.com/proffalken/monitorama-2024-demo-app
3. Create a new virtual env and run `pip install -r requirements.txt`
4. Execute `./bin/launch.sh`
5. Run the following: `cd basic_otel && ./manage.py createsuperuser`, then fill out the prompts
6. Log in to [http://localhost:8898/admin](http://localhost:8898/admin) and add a couple of spaces based on the [SpaceAPI Directory](https://github.com/SpaceApi/directory/blob/master/directory.json)
7. Run `./bin/launch_weather.sh`
8. Run `cd goclient && OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318" go run .`

You should start to see the traces and metrics flow through the platform and into Grafana Cloud App O11y. 
