# README

This script checks public GitHub releases with the help of GitHub's _releases_ page and reports back to upda. It won't
work on git tags!

To use, follow these instructions:

1. Copy the shell script to your PATH so you can use it, you also need the `upda` binary on your PATH.
2. Adapt the `REPOS=` in the script. Make sure the syntax follows GitHub's spaces format, so owner followed by a slash
   and then the repository's name.

To use the script with systemd, you need a `.service` file. Make sure to adapt `Environment` and `EnvironmentFile`. You
should not have secrets, like upda's Webhook ID or its token inside the environment directly. Sourcing it from a file is
better.

```ini
[Unit]
Description=Check public releases on GitHub and report to upda

[Service]
Environment="UPDA_SERVER_URL=..."
DynamicUser=true
EnvironmentFile=<path/to/secret-file/for/webhook-id-and-token/github-public-release-checker-env
ExecStart=github-public-release-checker.sh
Group=ghprc
StateDirectory=ghprc
Type=oneshot
User=ghprc
WorkingDirectory=/var/lib/ghprc
```

In addition, you need a `.timer` file for regularly running the service.

```ini
[Unit]
Description=Regularly check public releases on GitHub and report to upda

[Timer]
OnCalendar=weekly
Persistent=true
RandomizedDelaySec=1h

[Install]
WantedBy=timers.target
```

Example output from `journalctl`:

```shell
github-public-release-checker: [mailcow/mailcow-dockerized] New release detected: 2025-09a
github-public-release-checker: [ngosang/restic-exporter] New release detected: 1.7.0
github-public-release-checker: [promhippie/hcloud_exporter] New release detected: v3.2.7
github-public-release-checker: [promhippie/hetzner_exporter] New release detected: v2.2.6
github-public-release-checker: [crazy-max/diun] New release detected: v4.30.0
```