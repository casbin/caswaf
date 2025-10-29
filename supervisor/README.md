# Auto-Recovery Mechanism

## Overview

CasWAF includes a built-in self-recovery mechanism that automatically restarts the process when it crashes. This feature ensures high availability and minimizes downtime.

## Features

- **Cross-platform**: Works on Windows, Linux, macOS, and other platforms
- **Pure Go**: No cgo dependencies, fully portable
- **Configurable**: Can be enabled/disabled via configuration
- **Smart restart logic**: Prevents restart loops with configurable limits
- **Signal forwarding**: Properly handles shutdown signals

## How It Works

The auto-recovery mechanism uses a supervisor pattern:

1. When CasWAF starts with `autoRecoveryEnabled = true`, it spawns a supervisor process
2. The supervisor monitors the main CasWAF process
3. If the main process crashes, the supervisor automatically restarts it
4. If the process crashes repeatedly (5 times within 5 minutes by default), the supervisor stops to prevent infinite restart loops
5. On clean shutdown (SIGINT/SIGTERM), the supervisor forwards the signal and exits gracefully

## Configuration

To enable auto-recovery, add the following to your `conf/app.conf`:

```ini
autoRecoveryEnabled = true
```

To disable it:

```ini
autoRecoveryEnabled = false
```

## Default Settings

- **Maximum restarts**: 5 restarts within the restart window
- **Restart window**: 5 minutes
- **Restart delay**: 2 seconds between restarts

These settings are defined in `supervisor/supervisor.go` and can be customized if needed.

## Example

When auto-recovery is enabled and the process crashes:

```
Starting CasWAF with auto-recovery mechanism...
Started supervised process with PID: 12345
Process crashed after 2m30s: exit status 1
Waiting 2s before restarting... (restart 1/5)
Started supervised process with PID: 12346
```

## Technical Details

The supervisor uses standard Go libraries:
- `os/exec`: For spawning and managing child processes
- `os/signal`: For handling shutdown signals
- `syscall`: For signal forwarding

The implementation is platform-independent and uses only pure Go code, ensuring maximum portability.

## When to Use

Enable auto-recovery when:
- Running CasWAF in production environments
- You need automatic recovery from unexpected crashes
- You want to minimize manual intervention

Disable auto-recovery when:
- Debugging the application (to see the actual crash)
- Running under external process managers (like systemd, supervisord, or Docker restart policies)
- Testing and development

## Interaction with External Process Managers

If you're using an external process manager (systemd, supervisord, Docker, Kubernetes), you should typically disable the built-in auto-recovery to avoid conflicts:

- **systemd**: Use systemd's `Restart=` directive
- **Docker**: Use Docker's `--restart` policy
- **Kubernetes**: Use pod restart policies
- **supervisord**: Use supervisord's `autorestart` setting

In these cases, set `autoRecoveryEnabled = false` to let the external manager handle restarts.
