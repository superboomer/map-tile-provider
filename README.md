# **WARNING**

Make sure that you agree and stick to the policies
of the tile-providers before downloading!

[Tile usage policy](https://wiki.openstreetmap.org/index.php/Tile_usage_policy)
of OpenStreetMap if you use OSP tiles!

---
<div align="center">
  <img class="logo" src="https://raw.githubusercontent.com/superboomer/map-tile-provider/master/assets/logo.png" width="128px" height="128px" alt="logo"/>
  <br>
  <br>
  <b>Map Tile Provider</b>
  <br>
  <br>

  [![build](https://github.com/superboomer/map-tile-provider/actions/workflows/build.yml/badge.svg)](https://github.com/superboomer/map-tile-provider/actions/workflows/build.yml)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/superboomer/map-tile-provider)](https://goreportcard.com/report/github.com/superboomer/map-tile-provider)&nbsp;[![Coverage Status](https://coveralls.io/repos/github/superboomer/map-tile-provider/badge.svg?branch=master)](https://coveralls.io/github/superboomer/map-tile-provider?branch=master)
</div>


---
#### Environment variables

| Name          | Description   |  Optional | Default | 
| ------------- |:-------------:|:--------:| ------ |
|  ***LOG*** |
| LOG_SAVE  | enable logs save | ***Optional***  | false
| LOG_PATH     | logs path      | ***Optional***  | ./data/logs/log.jsonl
| LOG_MAX_BACKUPS | max backups count      |  ***Optional***  | 3
| LOG_MAX_SIZE | max logs size in megabytes      |  ***Optional***  | 1
| LOG_MAX_AGE | max logs age      |  ***Optional***  | 7
|  ***CACHE*** |
| CACHE_ENABLE | enable tile cache     | ***Optional***  | false
| CACHE_PATH | a path for cache directory     | ***Optional***  | ./data/cache
| CACHE_ALIVE | cache alive in minutes     | ***Optional***  | 14400
|  ***OTHERS*** |
| SCHEMA | providers specs    |  ***Required***  | *NO_DEFAULT*
| API_PORT | api port    |  ***Optional***  | 8080
| SWAGGER | swagger docs    |  ***Optional***  | false
| MAX_SIDE | max square side    |  ***Optional***  | 10
> All environment variables are available in [source code](https://github.com/superboomer/map-tile-provider/blob/master/app/options/opt.go)
***


# **Providers**

Example [providers.json](https://github.com/superboomer/map-tile-provider/blob/master/example/providers.json) contains 3 providers. *(but you can set up any providers as you wish. also service support loading .json from local FS)*

- OpenStreetMap
- Google Maps (Satellite)
- ArcGIS (Satellite)

> Don't forget about providers ToS

# **Docker Deploy**

You can easly deploy it via docker. Basic ***docker-compose.yml*** may look like this:
```YAML
version: '3.7'

services:

  map-tile-provider:
    image: ghcr.io/superboomer/map-tile-provider:latest
    container_name: map-tile-provider
    restart: unless-stopped
    environment:
      - SCHEMA=https://raw.githubusercontent.com/superboomer/map-tile-provider/master/example/providers.json
    ports:
      - "8080:8080"
```
> Full example [here](https://github.com/superboomer/map-tile-provider/blob/master/example)

***

