# **WARNING**

Make sure that you agree and stick to the policies
of the tile-providers before downloading!

[Tile usage policy](https://wiki.openstreetmap.org/index.php/Tile_usage_policy)
of OpenStreetMap if you use OSP tiles!

---
# ![alt text](assets/icon.png "Logo") **Map Tile Provider**

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
| API_PORT | api port    |  ***Optional***  | 8080
| SWAGGER | swagger docs    |  ***Optional***  | false
> All environment variables are available in [source code](https://github.com/mightrl/media-storage-service/blob/master/app/options/opt.go)
***


# **Providers**

- OpenStreetMap
- Google Maps (satellite)
- ArcGIS (satellite)

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
    ports:
      - "8080:8080"
```
***

