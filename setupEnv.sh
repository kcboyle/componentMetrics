#!/usr/bin/env bash

cf set-env componentMetrics DOPPLER_ADDR "ws://doppler.bosh-lite.com"
cf set-env componentMetrics CF_ACCESS_TOKEN "$(cf oauth-token | grep bearer)"

cf restage componentMetrics
