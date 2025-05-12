# deCONZ HomeKit Bridge

A HomeKit bridge for deCONZ/Phoscon that allows you to control your Zigbee devices through Apple HomeKit.

## Description

This project creates a bridge between deCONZ (a Zigbee gateway) and Apple HomeKit, allowing you to control your Zigbee devices using Apple's Home app, Siri, and other HomeKit-compatible applications. It supports a wide range of Zigbee devices including lights, sensors, and buttons.

## Requirements

* A running deCONZ/Phoscon gateway
* Docker (for containerized deployment)

## Installation

This application is designed to be installed and run via Docker only.

### Using Docker

1. Pull the Docker image:

   ```bash
   docker pull ghcr.io/0x2321/deconz-homekit:beta
   ```

2. Run the container:

   ```bash
   docker run -d -e DECONZ_IP=<your-deconz-ip> -e DECONZ_PORT=80 -v ./data:/data ghcr.io/0x2321/deconz-homekit:beta
   ```

### Using Docker Compose

Create a `docker-compose.yaml` file with the following content:

```yaml
version: '3'

services:
  deconz-homekit:
    image: ghcr.io/0x2321/deconz-homekit:beta
    container_name: deconz-homekit
    restart: unless-stopped
    environment:
      - DECONZ_IP=<your-deconz-ip>
      - DECONZ_PORT=80
    volumes:
      - ./data:/data
    network_mode: host
```

Then start the container with:

```bash
docker-compose up -d
```

> **Note:** The `./data` volume mount is used to persist the database file which stores the configuration and HomeKit pairing information. Make sure this directory exists and is writable.

## Configuration

The application requires the following environment variables:

* `DECONZ_IP`: The IP address of your deCONZ gateway
* `DECONZ_PORT`: The port of your deCONZ gateway (default 80)

On first run, the application will attempt to obtain an API key from your deCONZ gateway. You'll need to press the "Link" button on your deCONZ gateway when prompted.

## Usage

1. Start the application using one of the installation methods above.
2. The application will generate a HomeKit pairing code, which will be displayed in the logs.
3. Open the Home app on your iOS device, tap "Add Accessory", and enter the pairing code.
4. Once paired, your Zigbee devices will appear in the Home app and can be controlled through HomeKit.

## Device Support

This bridge creates a seamless connection between your deCONZ/Phoscon Zigbee devices and Apple HomeKit. It translates the capabilities of your Zigbee devices into HomeKit accessories, allowing you to control them through the Apple Home app, Siri voice commands, and HomeKit automations.

### Implementation Status

The tables below provide a comprehensive overview of all device types that can potentially be supported by this bridge. Each entry includes:

- **Description**: A user-friendly explanation of the device type and its functionality
- **deCONZ Type**: The internal identifier used by deCONZ to classify the device
- **Implemented**: Current implementation status (✓ = fully implemented, ❌ = not yet implemented)

Devices marked as implemented (✓) are fully functional in HomeKit and can be controlled through the Apple Home app. Devices not yet implemented (❌) are recognized by deCONZ but not currently exposed to HomeKit by this bridge.

If you have a specific device that's not yet implemented, consider contributing to the project or opening an issue on GitHub.

### Sensors

The following table lists all **device sensor** classes:

| Description              | deCONZ Type       | Implemented |
|:-------------------------|:------------------|:-----------:|
| Air Quality Sensor       | ZHAAirQuality     |      ❌      |
| Alarm Sensor             | ZHAAlarm          |      ❌      |
| Carbon Monoxide Sensor   | ZHACarbonMonoxide |      ❌      |
| Consumption Meter        | ZHAConsumption    |      ❌      |
| Fire Sensor              | ZHAFire           |      ❌      |
| Humidity Sensor          | ZHAHumidity       |      ❌      |
| Light Level Sensor       | ZHALightLevel     |      ❌      |
| Open/Close Sensor        | ZHAOpenClose      |      ✓      |
| Power Sensor             | ZHAPower          |      ❌      |
| Presence (Motion) Sensor | ZHAPresence       |      ✓      |
| Switch                   | ZHASwitch         |      ✓      |
| Pressure Sensor          | ZHAPressure       |      ❌      |
| Temperature Sensor       | ZHATemperature    |      ❌      |
| Time Sensor              | ZHATime           |      ❌      |
| Thermostat               | ZHAThermostat     |      ❌      |
| Vibration Sensor         | ZHAVibration      |      ❌      |
| Water Leak Sensor        | ZHAWater          |      ✓      |

### Lights

The following table lists all **light** classes:

| Description                                        | deCONZ Type              | Implemented |
|:---------------------------------------------------|:-------------------------|:-----------:|
| Basic light that can only be turned on or off      | On/Off Light             |      ✓      |
| Light with brightness control                      | Dimmable Light           |      ✓      |
| Light with adjustable white color temperature      | Color Temperature Light  |      ✓      |
| Light with RGB color control                       | Color Light              |      ❌      |
| Light with RGB and white color temperature control | Extended Color Light     |      ❌      |
| Smart plug that can be turned on or off            | On/Off Plug-in Unit      |      ✓      |
| Smart plug with power/brightness control           | Dimmable Plug-in Unit    |      ✓      |

## Development

For development, you can use the watch mode to automatically rebuild and run the application when changes are detected:

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

* [brutella/hap](https://github.com/brutella/hap) - HomeKit Accessory Protocol implementation in Go
* [deCONZ REST API](https://github.com/dresden-elektronik/deconz-rest-plugin) - The REST API for deCONZ
