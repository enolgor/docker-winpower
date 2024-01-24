# WINPOWER DOCKER

Simple docker image to use powerwalker winpower agent with UPS that support it.

Sample docker command:

`docker run --name winpower -d --privileged -p 8888:8888 -v "/dev/bus/usb:/dev/bus/usb" enolgor/winpower:latest`

Sample docker-compose:

```
version: "3.7"
services:
  winpower:
    image: enolgor/winpower:latest
    privileged: true
    ports:
      - 8888:8888
    volumes:
      - /dev/bus/usb:/dev/bus/usb
    restart: unless-stopped
```

You can then access the web interface at `https://localhost:8888` or poll the state of the UPS at `https://localhost:8888/0/json` in order to run further automations.

For example, [upsmon](https://github.com/enolgor/docker-winpower/tree/main/upsmon) is a service that fetches the status from winpower, can post status changes to a web service and run a script after a timeout on "AC Fail" status.