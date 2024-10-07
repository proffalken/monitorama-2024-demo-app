# Basic OTEL Demo

Runs a couple of microservices and reports to a local OTEL Collector/Alloy installation

## Get started

   1. Create a MySQL database: `CREATE DATABASE otel_basic;`
   2. Create a MySQL user: `CREATE USER otel_basic IDENTIFIED BY 'otel_basic';`
   3. Grant the required permissions: `GRANT ALL PRIVILEGES ON otel_basic.* TO otel_basic;FLUSH PRIVILEGES;`
   4. Install the requirements: `python -m env .venv;source .venv/bin/activate;pip install -r requirements.txt`
   5. Run the initial migrations and create a superuser: `cd basic_otel;./manage.py migrate;./manage.py createsuperuser;cd ..`
   6. Start the server: `./bin/launch.sh`
   7. Browse to [http://localhost:8898/admin](http://localhost:8898/admin) and log in using your superuser account created in step (4)
   8. Add a new "space" to the platform, "Make Monmouth" as the name and "https://members.makemonmouth.co.uk/api/spacedirectory/" will work as a test but any of the values from [this file](https://github.com/SpaceApi/directory/blob/master/directory.json_) will work if you want to find one closer to you
   9. Launch the weather service (this does not need any migrations to be run): `./bin/launch_weather.sh`
   10. Run the GO client to poll the API every 2 seconds: `cd goclient;OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318/" go run .`

As long as you have followed the above and have a local running OTEL collector, you should start to see a service map appear in your O11y tooling.
