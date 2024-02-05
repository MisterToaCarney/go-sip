# A toy sip implementation in go

You need a config file `config.go` in package root in order to build. Below is an example.

```{go}
package main

import "github.com/MisterToaCarney/gosip/utils"

func GetConfig() utils.MyDetails {
  return utils.MyDetails{
    Username:   "0800yourphone",
    Password:   "yourpassword",
    RemoteHost: "remote sip trunk",
    RemotePort: "5060",
    Scheme:     "sip",
    Transport:  "tcp",
  }
}
```
