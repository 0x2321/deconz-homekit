> **English Version:** [README.en.md](README.en.md)

# deCONZ HomeKit Bridge

Eine HomeKit-Bridge für deCONZ/Phoscon, die es ermöglicht, deine Zigbee-Geräte über Apple HomeKit zu steuern.

## Installation

Diese Anwendung ist dafür ausgelegt, ausschließlich über Docker installiert und betrieben zu werden.

### Verwendung mit Docker

1. Ziehe das Docker-Image:

```bash
docker pull ghcr.io/0x2321/deconz-homekit:main
```

2. Starte den Container:

```bash
docker run -d \
  -e DECONZ_IP=<deconz-ip> \
  -e DECONZ_PORT=80 \
  -v ./data:/data \
  ghcr.io/0x2321/deconz-homekit:main
```

### Verwendung mit Docker Compose

Erstelle eine `docker-compose.yaml` mit folgendem Inhalt:

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

Starte dann den Container mit:

```bash
docker-compose up -d
```

> **Hinweis:** Der `./data`-Volume-Mount dient dazu, die Datenbankdatei zu speichern, in der die Konfiguration und die HomeKit-Pairing-Informationen abgelegt werden. Stelle sicher, dass dieses Verzeichnis existiert und beschreibbar ist.

## Konfiguration

Stelle folgende Umgebungsvariablen ein:

* `DECONZ_IP`: IP-Adresse des deCONZ-Gateways
* `DECONZ_PORT`: Port des deCONZ-Gateways (Standard: 80)

Beim ersten Start fordert die Anwendung einen API-Key vom Gateway an. Öffne dazu die Phoscon Web App, navigiere zu **Einstellungen → Gateway → Erweiterte Einstellungen** und klicke auf **"App authentifizieren"**, um den Zugriff zu autorisieren.

## Geräteunterstützung

Nicht alle deCONZ-Gerätekategorien sind aktuell implementiert.

Statussymbole: ✅ umgesetzt | 🧪️ in Test (Feedback erwünscht) | ❌ ausstehend

#### Sensoren

| Beschreibung             | deCONZ Typ        | Implementiert |
| ------------------------ | ----------------- | ------------- |
| Öffnungs-/Schließsensor  | ZHAOpenClose      | ✅             |
| Präsenz-/Bewegungssensor | ZHAPresence       | ✅             |
| Schalter                 | ZHASwitch         | ✅             |
| Wasserlecksensor         | ZHAWater          | 🧪             |
| Luftgütesensor           | ZHAAirQuality     | ❌             |
| Alarmsensor              | ZHAAlarm          | ❌             |
| Kohlenmonoxid-Sensor     | ZHACarbonMonoxide | ❌             |
| Verbrauchszähler         | ZHAConsumption    | ❌             |
| Feuermelder              | ZHAFire           | ❌             |
| Feuchtigkeitssensor      | ZHAHumidity       | ❌             |
| Lichtsensor              | ZHALightLevel     | ❌             |
| Leistungssensor          | ZHAPower          | ❌             |
| Drucksensor              | ZHAPressure       | ❌             |
| Temperatursensor         | ZHATemperature    | ❌             |
| Zeitsensor               | ZHATime           | ❌             |
| Thermostat               | ZHAThermostat     | ❌             |
| Vibrationssensor         | ZHAVibration      | ❌             |

#### Lichter

| Gerätekategorie                                | deCONZ Typ              | Status |
|------------------------------------------------|-------------------------|--------|
| Einfaches Licht (nur Ein/Aus)                  | On/Off Light            | ✅      |
| Licht mit Helligkeitssteuerung                 | Dimmable Light          | ✅      |
| Licht mit einstellbarer Weißfarbtemperatur     | Color Temperature Light | ✅      |
| Intelligente Steckdose (Ein/Aus)               | On/Off Plug-in Unit     | ✅      |
| Intelligente Steckdose mit Dimmfunktion        | Dimmable Plug-in Unit   | ✅      |
| Licht mit RGB-Farbsteuerung                    | Color Light             | ❌      |
| Licht mit RGB- und Weißfarbtemperatursteuerung | Extended Color Light    | ❌      |

## Entwicklung

Für die Entwicklung kannst du den Watch-Mode verwenden, um die Anwendung bei Änderungen automatisch neu zu bauen und zu starten:

```bash
make watch
```

## Lizenz

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

## Danksagungen

* [brutella/hap](https://github.com/brutella/hap) - HomeKit Accessory Protocol implementation in Go
* [deCONZ REST API](https://github.com/dresden-elektronik/deconz-rest-plugin) - The REST API for deCONZ
