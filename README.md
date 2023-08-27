## Main

This app is a simple CLI tool to invoke the service menu on LG televisions.
I got extremely fed up with my television dimming the whole screen, despite the
normal dimming settings being turned off, when watching videos that have very
little movement for the majority of the picture (e.g. some YouTube videos). I
don't have a system to use other tools with, so I spent an afternoon hacking
this one together.

The code is pieced together by reading how the following projects do things:

+ https://github.com/Maassoft/ColorControl
+ https://github.com/Danovadia/lgtv-http-server
+ https://github.com/supersaiyanmode/PyWebOSTV
+ https://github.com/merdok/homebridge-webos-tv/tree/master

## Building

```sh
$ go build
```

## Using

1. Run the app: `$ ./lgtvtool 10.0.0.15` (where `10.0.0.15` is the IP address
of the television). 
2. Enter passcode with television remote number buttons. `0413` worked on my
television. A list of passcodes can be found on
[https://www.wikihow.com/Display-the-Secret-Menu-in-LG-TVs](https://www.wikihow.com/Display-the-Secret-Menu-in-LG-TVs).
3. Navigate to the "OLED" menu.
4. Disable the "GSR" setting. (
https://www.reddit.com/r/OLED/comments/x42to5/lg_c2_share_your_ways_to_limit_abl/imvc8hn/
 indicates a "TCL" setting should also be disabled, but I don't see it on my
television)
5. Press the "gear" button on the remote to close the menu.
6. Use `ctrl+c` to terminate the tool.
