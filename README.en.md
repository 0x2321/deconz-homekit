# deCONZ HomeKit Bridge

A HomeKit bridge for deCONZ/Phoscon that enables you to control your Zigbee devices via Apple HomeKit.

## Installation

This application is designed to be installed and operated exclusively via Docker.

### Using Docker

1. Pull the Docker image:

```bash
docker pull ghcr.io/0x2321/deconz-homekit:main
```

2. Start the container:

```bash
docker run -d \
  -e DECONZ_IP=<deconz-ip> \
  -e DECONZ_PORT=80 \
  -v ./data:/data \
  ghcr.io/0x2321/deconz-homekit:main
```

### Using Docker Compose

Create a `docker-compose.yaml` with the following content:

```yaml
services:
   deconz-homekit:
      image: ghcr.io/0x2321/deconz-homekit:main
      container_name: deconz-homekit
      restart: unless-stopped
      environment:
         - DECONZ_IP=<deconz-ip>
         - DECONZ_PORT=80
      volumes:
         - ./data:/data
      network_mode: host
```

Then start the container with:

```bash
docker-compose up -d
```

> **Note:** The `./data` volume mount is used to persist the database file that stores the configuration and HomeKit pairing information. Ensure this directory exists and is writable.

## Configuration

Set the following environment variables:

* `DECONZ_IP`: IP address of the deCONZ gateway
* `DECONZ_PORT`: Port of the deCONZ gateway (default: 80)

On the first start, the application will request an API key from the gateway. To authorize access, open the Phoscon web app, navigate to **Settings → Gateway → Advanced Settings**, and click **“Authenticate app”**.

## Device Support

Not all deCONZ device categories are currently implemented.

Status icons: ✅ implemented | 🧪️ in testing (feedback welcome) | ❌ pending

#### Sensors

| Description             | deCONZ Type       | Implemented |
| ----------------------- | ----------------- | ----------- |
| Open/Close sensor       | ZHAOpenClose      | ✅           |
| Presence/Motion sensor  | ZHAPresence       | ✅           |
| Switch                  | ZHASwitch         | ✅           |
| Water leak sensor       | ZHAWater          | 🧪           |
| Air quality sensor      | ZHAAirQuality     | ❌           |
| Alarm sensor            | ZHAAlarm          | ❌           |
| Carbon monoxide sensor  | ZHACarbonMonoxide | ❌           |
| Power consumption meter | ZHAConsumption    | ❌           |
| Smoke detector          | ZHAFire           | ❌           |
| Humidity sensor         | ZHAHumidity       | ❌           |
| Light level sensor      | ZHALightLevel     | ❌           |
| Power sensor            | ZHAPower          | ❌           |
| Pressure sensor         | ZHAPressure       | ❌           |
| Temperature sensor      | ZHATemperature    | ❌           |
| Time sensor             | ZHATime           | ❌           |
| Thermostat              | ZHAThermostat     | ❌           |
| Vibration sensor        | ZHAVibration      | ❌           |

#### Lights

| Device Category                                 | deCONZ Type             | Status |
| ----------------------------------------------- | ----------------------- | ------ |
| Basic light (on/off only)                       | On/Off Light            | ✅      |
| Light with brightness control                   | Dimmable Light          | ✅      |
| Light with adjustable white color temperature   | Color Temperature Light | ✅      |
| Smart plug (on/off)                             | On/Off Plug-in Unit     | ✅      |
| Smart plug with dimming function                | Dimmable Plug-in Unit   | ✅      |
| Light with RGB color control                    | Color Light             | ❌      |
| Light with RGB and white color temperature ctrl | Extended Color Light    | ❌      |

## Development

For development, you can use the watch mode to automatically rebuild and restart the application upon changes:

```bash
make watch
```

## License

MIT License
Copyright (c) 2025 Bastian Dietrich

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Acknowledgements

* [brutella/hap](https://github.com/brutella/hap) – HomeKit Accessory Protocol implementation in Go
* [deCONZ REST API](https://github.com/dresden-elektronik/deconz-rest-plugin) – The REST API for deCONZ
