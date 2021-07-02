## Securum Exire Traefik Plugin

## Using the plugin

Register your traefik instance [here](https://pilot.traefik.io/instances?action=register).

After registering you will get the traefik pilot token.

##  Static Configuration

YAML:
```
pilot:
    token: "xxxx"
experimental:
    plugins:
        traefik-plugin-securum-exire:
            moduleName: "github.com/mayankkumar2/traefik-plugin-securum-exire"
            version: "0.1.5"
```     
TOML:
```
[pilot]
  token = "xxxx"
[experimental.plugins.traefik-plugin-securum-exire]
  moduleName = "github.com/mayankkumar2/traefik-plugin-securum-exire"
  version = "0.1.5"
```
CLI:
```
--pilot.token=xxxx
--experimental.plugins.traefik-plugin-securum-exire.modulename=github.com/mayankkumar2/traefik-plugin-securum-exire
--experimental.plugins.traefik-plugin-securum-exire.version=0.1.5
```  
## Dynamic configuration
(change the url with the configured leak server url)<br>
YAML:
```
http:
    middlewares:
        my-traefik-plugin-securum-exire:
            plugin:
                traefik-plugin-securum-exire:
                    url: "0.0.0.0:10012"
``` 
TOML:
```toml
[http]
  [http.middlewares]
    [http.middlewares.my-traefik-plugin-securum-exire]
      [http.middlewares.my-traefik-plugin-securum-exire.plugin]
        [http.middlewares.my-traefik-plugin-securum-exire.plugin.traefik-plugin-securum-exire]
          url = "0.0.0.0:10012"
```
