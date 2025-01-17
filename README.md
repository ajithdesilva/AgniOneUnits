# AgniOne Demo Unit.

Application Units for Agni Application Framework

This is a demo AgniOne Unit that will utilize the AgniOne HTTP plugin

This unit will be hosted & controlled  by the AgniOne Application Framework



Basically, to launch AgniOne Unit, it AgniOne Framework requires 3 files.
1. Application specific configuration (app.config)
1. AgniOne Unit binary  (httpdemo.so)
2. AgniOne Unit's specific configuration file (httpdemo.config)
   

1. app.config  -> this will define the properties that AgniOne framework need to use during the initialization.
   
   
   ```
   {
    "app": {
      "name": "AgniOne Demo Application",
      "id": "AgniOne_demo_http",
      "version": "1.0.0.0"
    },
    "log":{
      "leg_level": "debug",
      "log_file_max_size": 10000,
      "log_file_base_path": "/var/log/app/"
    },
    "appunits": [
      {
        "uname": "Http_demo",
        "enable": 1,
        "pool_size":1,
        "path":"<UNIT-BINARY-PATH>",
        "config":"<UNIT-CONFIG-PATH>"
      }
    ]
    }
   ```
   
   
2. httpdemo.config -> this file contains the configuration that needed by demohttp Unit to function as required.

```
    {
        "get":{
            "plugin_type":"default",
            "url":"https://httpbin.org/get",
            "frequancy_secs":3
        },
        "post":{
            "plugin_type":"default",
            "url":"https://httpbin.org/post",
            "frequancy_secs":3
        }
    }
    

```


### build

Using the build script Unit can be built & deployed easily.
Build script requires the path of the AgniOne Application framework to deploy.
    Eg:-
       ./build.sh <AgniOne-Framework-Path>

```
    ./build.sh ~/AgniOneFM/AgniOne
```

During the build process, it will build the binary of the HTTP demo unit and deploy it to the given Framework path.

Also, will update the unit's app.config with relevant paths, so that there is no need to modify the Unit path and it's config path
manually.



