# Clock In On Startup

> This example is only for GNU/Linux users w/ systemd.

Systemd is the standard service manager used in many GNU/Linux distributions for managing system and user services. With systemd, you can set up automatic clock-in on startup. Please keep in mind that `toffu` considers the following:

- Whether you have manually clocked in already.
- Whether today is a workday.
- Whether you have worked the scheduled hours.

## Install Toffu

Check out the [installation instructions](../../README.md#installation) first.

## Create the systemd config file

Copy the next [systemd config file](toffu-in.service) in `$HOME/.config/systemd/user/toffu-in.service`:

```ini
[Unit]
Description=Toffu in on startup
After=dbus.service network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/toffu in

[Install]
WantedBy=default.target
```

## Enable the service

Start the service:

```shell
systemctl --user start toffu-in.service
```

And check if everything is running fine:

```shell
systemctl --user status toffu-in.service
```

Once the service is working, you can enable it on startup:

```shell
systemctl --user enable toffu-in.service
```
