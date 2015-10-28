#!/usr/bin/env bash

cf set-env componentMetrics DOPPLER_ADDR "wss://doppler.oak.cf-app.com:4443"
cf set-env componentMetrics CF_ACCESS_TOKEN "$(cf oauth-token | grep bearer)"

cf restart componentMetrics
