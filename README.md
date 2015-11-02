# ComponentMetrics
To run this app, please have [Cloud Foundry](https://github.com/cloudfoundry/cf-release) installed.

ComponentMetrics need to have access to your doppler websocket (`wss://`) endpoint and your cf oauth token (`$(cf oauth-token | grep bearer)`). These need to be exported in variables named `DOPPLER_ADDR` and `CF_ACCESS_TOKEN`. Also export the `PORT` variable to any open port of your choosing.

To run the application, simply run `go run main.go` inside the repository. Then navigate to `localhost:(PORT)/messages` to see the graph. 

To deploy to Cloud Foundry, simply change the DOPPLER_ADDR inside of the `setupEnv.sh` script, and execute `deploy.sh`. This will deploy componentMetrics to Cloud Foundry and setup the environment variables the app needs to run. Then you can navigate to the address reserved for your app (the `/messages` endpoint) and view the metrics there.
