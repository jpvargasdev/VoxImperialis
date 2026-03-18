# Vox Imperialis

A deterministic, command-driven XMPP operations bot for homelab management.
It connects to a Prosody XMPP server as a client and responds to commands sent
by a configured list of authorised JIDs. No LLMs. No shell expansion. Secure by design.

---

## Requirements

| Dependency | Notes |
|---|---|
| Go 1.23+ | |
| XMPP server | Prosody recommended |
| `lm-sensors` | Optional — needed for the `sensors` command |
| `systemd` | Required for `service` commands |

---

## Setup

### 1. Clone and enter the directory

```sh
cd /opt/vox-imperialis
```

### 2. Configure

```sh
cp .env.example .env
$EDITOR .env
```

Set your XMPP credentials, allowed users, and the services the bot may manage.

### 3. Fetch dependencies

```sh
go mod tidy
```

### 4. Build

```sh
go build -o vox-imperialis .
```

### 5. Run

```sh
./vox-imperialis
```

---

## Docker deployment (with Prosody)

VoxImperialis is deployed alongside a Prosody XMPP server via Docker Compose,
with Traefik handling TCP routing for XMPP ports.

### Register XMPP users in Prosody

After the Prosody container is running, create the bot account and any
authorised user accounts:

```sh
# Create the bot account
podman exec prosody prosodyctl adduser vox-imperialis@vox.example.com
# You will be prompted for a password

# Create an authorised user account
podman exec prosody prosodyctl adduser operator@vox.example.com
```

Then set the matching credentials in the VoxImperialis environment:

```
XMPP_JID=vox-imperialis@vox.example.com
XMPP_PASSWORD=<password-set-above>
XMPP_SERVER=vox.example.com:5222
ALLOWED_USERS=operator@vox.example.com
ALLOWED_SERVICES=myservice
```

---

## Commands

| Command | Description |
|---|---|
| `help` | List available commands |
| `status` | Uptime, load average, memory usage, root disk usage |
| `sensors` | Hardware sensor readings from `lm-sensors` |
| `service status <name>` | Show systemd service status |
| `service start <name>` | Start a systemd service |
| `service stop <name>` | Stop a systemd service |
| `service restart <name>` | Restart a systemd service |

### Example session

```
you:  status
bot:  [status]
      host:   homelab
      uptime: 12d 4h 31m
      load:   0.24 0.31 0.28
      memory: 8.1G / 32.0G
      disk /:  41.3G / 118.0G (35%)

you:  service restart jellyfin
bot:  [service]
      name:   jellyfin
      action: restart
      result: success
```

---

## Run as a systemd service

```sh
# Create a dedicated user (optional but recommended)
sudo useradd -r -s /sbin/nologin vox-imperialis

# Install the binary and config
sudo mkdir -p /opt/vox-imperialis
sudo cp vox-imperialis /opt/vox-imperialis/
sudo cp .env          /opt/vox-imperialis/
sudo chown -R vox-imperialis:vox-imperialis /opt/vox-imperialis
sudo chmod 600 /opt/vox-imperialis/.env

# Install and enable the service
sudo cp vox-imperialis.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now vox-imperialis

# Follow logs
sudo journalctl -fu vox-imperialis
```

### Service management permissions

To allow the bot to call `systemctl start/stop/restart`, the service account
must have the necessary permissions. Using **polkit** is recommended over running
as root. Create `/etc/polkit-1/rules.d/50-vox-imperialis.rules`:

```js
polkit.addRule(function(action, subject) {
    if (action.id == "org.freedesktop.systemd1.manage-units" &&
        subject.user == "vox-imperialis") {
        return polkit.Result.YES;
    }
});
```

For a quick homelab setup you can also run the service as root by removing the
`User=` and `Group=` lines from the service file — understand the implications
before doing so.

---

## Security

- Only JIDs listed in `ALLOWED_USERS` can issue commands.
- Only services listed in `ALLOWED_SERVICES` can be managed.
- All system calls use `exec.Command` with explicit, discrete arguments — no shell
  expansion, no user-controlled command strings.
- TLS verification is **on** by default; set `XMPP_TLS_SKIP_VERIFY=true` only for
  self-signed homelab certificates.

---

## Project layout

```
.
├── main.go               entry point, wires components together
├── config.go             AppConfig, Load(), Get() — reads .env
├── auth.go               JID allowlist check
├── parser.go             text → Command struct
├── dispatcher.go         command name → HandlerFunc routing
├── format.go             error message helpers
├── xmpp_client.go        XMPP connection, reconnect loop, message routing
├── handlers/
│   ├── handlers.go       Command and HandlerFunc types
│   ├── help.go           help command
│   ├── status.go         status command
│   ├── sensors.go        sensors command
│   └── services.go       service command with allowlist enforcement
├── system/
│   ├── status.go         GetSystemStatus() via /proc
│   ├── sensors.go        GetSensors() via lm-sensors
│   └── systemd.go        SystemdStatus/Start/Stop/Restart via systemctl
├── .env.example          configuration template
└── vox-imperialis.service systemd unit file
```
