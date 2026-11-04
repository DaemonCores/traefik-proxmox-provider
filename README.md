# Traefik Proxmox Provider

<p align="center">
  <img src="https://raw.githubusercontent.com/DaemonCores/.github/refs/heads/main/assets/banner.svg" alt="DaemonCores Banner" width="100%"/>
</p>

<p align="center">
  <a href="https://github.com/orgs/DaemonCores/teams"><img src="https://img.shields.io/badge/Community-000000?style=for-the-badge&logoColor=white" alt="Community"/></a>
  <img src="https://img.shields.io/badge/Upstream-NX211_v0.8.1-A81D33?style=for-the-badge" alt="Upstream base"/>
</p>

A Traefik provider plugin that automatically configures routing based on Proxmox VE
virtual machines and containers — maintained by DaemonCores since 2026-11-04.

## Why DaemonCores maintains this fork

Upstream (NX211) stays active and we credit their work: the whole provider below v0.8.1
is theirs, under Apache-2.0. We took over maintenance for one precise reason: the
middleware-in-labels capability (upstream issue #39, PR #41) has been parked there
since February 2026, and DaemonCores production runs on the PVE traefik-provider
plugin — we could not keep waiting on that review calendar.

So this repository is now the maintained fork: bug fixes and features needed by
DaemonCores land here first. Anything broadly useful may be proposed back upstream,
but our release cadence no longer depends on it.

## What we changed (on top of upstream v0.8.1)

- **`parseMiddlewares` — middleware definitions parsed from CT/VM labels**
  (`traefik.http.middlewares.*` entries in instance metadata are merged into the
  dynamic configuration), which upstream's provider ignores.
- Case-insensitive label handling and TLS domains array syntax.
- Yaegi test target fix (hidden files + tmp directory handling).
- Dev container configuration, troubleshooting documentation.

## Installation

```yaml
experimental:
  plugins:
    traefik-proxmox-provider:
      moduleName: github.com/DaemonCores/traefik-proxmox-provider
      version: v0.9.0
```

(Upstream installs pointed at `github.com/NX211/traefik-proxmox-provider` — switch
the `moduleName` to this repository for the maintained line.)

## Configuration, usage, features

Unchanged from upstream v0.8.1 — see the original README sections kept below and
the upstream documentation. DaemonCores production notes:

- The PVE-side labels the provider reads are written as `traefik.*=value` lines in
  the guest `description` field, one per line, first line `traefik.enable=true`.
- A router referencing an unparsed middleware used to invalidate the whole payload;
  with `parseMiddlewares` the definitions travel with the guest that needs them.

## License

Apache-2.0 — inherited from upstream. DaemonCores contributions are published under
the same license. Full upstream history and authorship is preserved in git history on purpose.

## Maintenance

Maintained by [DaemonCores](https://github.com/DaemonCores) — issues and PRs go here.
Upstream (NX211) remains the origin project; thanks to all its contributors.