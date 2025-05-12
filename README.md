> **English Version:** [README.en.md](README.en.md)

# deCONZ HomeKit Bridge

Eine HomeKit-Bridge für deCONZ/Phoscon, die es ermöglicht, deine Zigbee-Geräte über Apple HomeKit zu steuern.

## Beschreibung

Dieses Projekt erstellt eine Bridge zwischen deCONZ (einem Zigbee-Gateway) und Apple HomeKit, sodass du deine Zigbee-Geräte mit der Apple Home App, Siri und anderen HomeKit-kompatiblen Anwendungen steuern kannst. Es werden zahlreiche Zigbee-Geräte unterstützt, einschließlich Lichter, Sensoren und Tasten.

## Anforderungen

* Ein laufendes deCONZ/Phoscon-Gateway
* Docker (für die containerisierte Bereitstellung)

## Installation

Diese Anwendung ist dafür ausgelegt, ausschließlich über Docker installiert und betrieben zu werden.

### Verwendung mit Docker

1. Ziehe das Docker-Image:

   ```bash
   docker pull ghcr.io/0x2321/deconz-homekit:main
   ```

2. Starte den Container:

   ```bash
   docker run -d -e DECONZ_IP=<deine-deconz-ip> -e DECONZ_PORT=80 -v ./data:/data ghcr.io/0x2321/deconz-homekit:main
   ```

### Verwendung mit Docker Compose

Erstelle eine `docker-compose.yaml` mit folgendem Inhalt:

```yaml
version: '3'

services:
  deconz-homekit:
    image: ghcr.io/0x2321/deconz-homekit:main
    container_name: deconz-homekit
    restart: unless-stopped
    environment:
      - DECONZ_IP=<deine-deconz-ip>
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

Die Anwendung benötigt folgende Umgebungsvariablen:

* `DECONZ_IP`: Die IP-Adresse deines deCONZ-Gateways
* `DECONZ_PORT`: Der Port deines deCONZ-Gateways (Standard: 80)

Beim ersten Start versucht die Anwendung, einen API-Key von deinem deCONZ-Gateway zu erhalten. Drücke dazu bei Aufforderung die "Link"-Taste auf deinem deCONZ-Gateway.

## Nutzung

1. Starte die Anwendung mit einer der oben beschriebenen Installationsmethoden.
2. Die Anwendung generiert einen HomeKit-Pairing-Code, der in den Logs angezeigt wird.
3. Öffne die Home App auf deinem iOS-Gerät, tippe auf "Zubehör hinzufügen" und gib den Pairing-Code ein.
4. Nach erfolgreichem Pairing erscheinen deine Zigbee-Geräte in der Home App und können über HomeKit gesteuert werden.

## Geräteunterstützung

Diese Bridge stellt eine nahtlose Verbindung zwischen deinen deCONZ/Phoscon Zigbee-Geräten und Apple HomeKit her. Sie übersetzt die Fähigkeiten deiner Zigbee-Geräte in HomeKit-Zubehörteile, sodass du sie über die Apple Home App, Siri-Sprachbefehle und HomeKit-Automationen steuern kannst.

### Implementierungsstatus

Die folgenden Tabellen bieten einen umfassenden Überblick aller Gerätetypen, die potenziell von dieser Bridge unterstützt werden können. Jeder Eintrag enthält:

* **Beschreibung**: Eine Erklärung des Gerätetyps und seiner Funktion
* **deCONZ Type**: Der interne Bezeichner, den deCONZ zur Klassifizierung des Geräts verwendet
* **Implemented**: Aktueller Implementierungsstatus (✓ = vollständig implementiert, ❌ = noch nicht implementiert)

Geräte mit dem Status Implemented (✓) sind in HomeKit voll funktionsfähig und können über die Apple Home App gesteuert werden. Geräte, die noch nicht implementiert sind (❌), werden von deCONZ erkannt, aber von dieser Bridge derzeit nicht an HomeKit weitergegeben.

#### Sensoren

Die folgende Tabelle listet alle **Device Sensor**-Klassen auf:

| Beschreibung             | deCONZ Type       | Implemented |
| :----------------------- | :---------------- | :---------: |
| Luftgütesensor           | ZHAAirQuality     |      ❌      |
| Alarmsensor              | ZHAAlarm          |      ❌      |
| Kohlenmonoxid-Sensor     | ZHACarbonMonoxide |      ❌      |
| Verbrauchszähler         | ZHAConsumption    |      ❌      |
| Feuermelder              | ZHAFire           |      ❌      |
| Feuchtigkeitssensor      | ZHAHumidity       |      ❌      |
| Lichtsensor              | ZHALightLevel     |      ❌      |
| Öffnungs-/Schließsensor  | ZHAOpenClose      |      ✓      |
| Leistungssensor          | ZHAPower          |      ❌      |
| Präsenz-/Bewegungssensor | ZHAPresence       |      ✓      |
| Schalter                 | ZHASwitch         |      ✓      |
| Drucksensor              | ZHAPressure       |      ❌      |
| Temperatursensor         | ZHATemperature    |      ❌      |
| Zeitsensor               | ZHATime           |      ❌      |
| Thermostat               | ZHAThermostat     |      ❌      |
| Vibrationssensor         | ZHAVibration      |      ❌      |
| Wasserlecksensor         | ZHAWater          |      ✓      |

#### Lichter

Die folgende Tabelle listet alle **Light**-Klassen auf:

| Beschreibung                                   | deCONZ Type             | Implemented |
| :--------------------------------------------- | :---------------------- | :---------: |
| Einfaches Licht (nur Ein/Aus)                  | On/Off Light            |      ✓      |
| Licht mit Helligkeitssteuerung                 | Dimmable Light          |      ✓      |
| Licht mit einstellbarer Weißfarbtemperatur     | Color Temperature Light |      ✓      |
| Licht mit RGB-Farbsteuerung                    | Color Light             |      ❌      |
| Licht mit RGB- und Weißfarbtemperatursteuerung | Extended Color Light    |      ❌      |
| Intelligente Steckdose (Ein/Aus)               | On/Off Plug-in Unit     |      ✓      |
| Intelligente Steckdose mit Dimmfunktion        | Dimmable Plug-in Unit   |      ✓      |

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
