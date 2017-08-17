## PMC - Poor Man's Check
PMC is a cheap tool to verify and register network endpoints.

The motivation for such tool was the many downlinks that I had with my cable company, and with the recorded info I can argue more about their quality.

### Config
The app looks for a file named pmc.conf written in toml format.

Example:

```
[[verify]]
type="http"
label="Awesome Identifier"
config="http://www.amazon.com"
every="5m"
timeout="1s"

[[verify]]
type="dns"
label="Some lookup check identifier"
config="google.com"
every="1d"
timeout="30ms"

[[verify]]
type="ping"
label="Raw socket ping identifier"
config="apple.com"
every="5m"
timeout="1s"

[[register]]
type="text"

[[register]]
type="influxdb"
config="http://username:password@host:port?db=dbname&series=seriesname"
```
Every verifier will write to all the registers.
Currently the example shows all the registers and verifiers available.

When using the ping verifier root privileges or CAP_NET_RAW on Linux

```
# setcap cap_net_raw,cap_net_admin=eip PATH_TO_PMC
```

If you don't want the app to output to the terminal just remove the text register.

### Linux service

#### systemd

Write a file at /lib/systemd/system/pmc.service with the following contents:

```
[Unit]
Description=PMC Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/pmc -conf=/etc/pmc.conf
User=root

[Install]
WantedBy=multi-user.target
```
Then execute:

``` 
systemctl daemon-reload 
systemctl start pmc
systemctl status pmc
```


More info about systemctl at https://www.digitalocean.com/community/tutorials/how-to-use-systemctl-to-manage-systemd-services-and-units

### TODO:

- systemd init scripts
- plotting capabilities
	
###Know Bugs

Sometimes the influxdb register complains with ```unsupported point type: time.Duration```
