# Traefik Proxmox Provider

<p align="center">
  <img src="https://raw.githubusercontent.com/DaemonCores/.github/refs/heads/main/assets/banner.svg" alt="DaemonCores Banner" width="100%"/>
</p>

<p>
  <strong align="left">Simplify and Innovate for Everyone.</strong>
  <a href="https://github.com/DaemonCores/traefik-proxmox-provider/wiki"><img align="right" src="https://img.shields.io/badge/Wiki-FFFFFF?style=for-the-badge&logoColor=white" alt="Documentation"/></a>
  <a href="https://github.com/orgs/DaemonCores/discussions"><img align="right" src="https://img.shields.io/badge/Community-000000?style=for-the-badge&logoColor=white" alt="Community"/></a>
  <a href="https://github.com/DaemonCores/traefik-proxmox-provider"><img align="right" src="https://img.shields.io/badge/Traefik_Proxmox_Provider-A81D33?style=for-the-badge&logo=traefikproxy&logoColor=white" alt="Traefik Proxmox Provider"/></a>
  
  <em>Identify gaps and fill them, make improvements where possible, but above all, empower developers to offer more to users.</em>
</p>

---

A Traefik provider that automatically configures routing based on Proxmox VE virtual machines and containers — discovered from each VM/container's notes field, no files to mount, no sidecar to run.

- Automatically discovers Proxmox VE virtual machines and containers
- Configures routing based on VM/container metadata
- Middleware definitions travel with the instance (`parseMiddlewares` — DaemonCores addition)
- HTTP and HTTPS endpoints, configurable polling interval, SSL validation options, full Traefik routing/middleware/TLS support

## Maintained by DaemonCores

Upstream ([NX211/traefik-proxmox-provider](https://github.com/NX211/traefik-proxmox-provider), v0.8.1 at takeover) stays credited: everything below v0.8.1 is their work under Apache-2.0 (full history preserved). Upstream's middleware-in-labels proposal (PR #41) has been parked since February 2026; DaemonCores production runs this plugin on its PVE edge, so the line is maintained here since 2026-11-04. Bug fixes and features needed by DaemonCores land in this repository first; useful patches may be proposed upstream without our releases depending on it.

The deltas on top of upstream v0.8.1:

- **`parseMiddlewares`** — `traefik.http.middlewares.*` definitions found in instance labels are merged into the dynamic configuration (upstream only parses routers/services/TLS).
- Case-insensitive label handling and TLS domains array syntax.
- Yaegi test target fix (hidden files and tmp directory handling), dev container, troubleshooting docs.

## Documentation

The full documentation lives in the [wiki](https://github.com/DaemonCores/traefik-proxmox-provider/wiki):

- [Installation](https://github.com/DaemonCores/traefik-proxmox-provider/wiki/Installation)
- [Configuration](https://github.com/DaemonCores/traefik-proxmox-provider/wiki/Configuration)
- [Proxmox API token setup](https://github.com/DaemonCores/traefik-proxmox-provider/wiki/PVE-API-Tokens)
- [VM/Container labeling](https://github.com/DaemonCores/traefik-proxmox-provider/wiki/VM-Container-Labeling)
- [Examples](https://github.com/DaemonCores/traefik-proxmox-provider/wiki/Examples)
- [How it works](https://github.com/DaemonCores/traefik-proxmox-provider/wiki/How-It-Works)
- [Troubleshooting](https://github.com/DaemonCores/traefik-proxmox-provider/wiki/Troubleshooting)
- [DaemonCores takeover](https://github.com/DaemonCores/traefik-proxmox-provider/wiki/DaemonCores-Takeover)

## Contributing

Contributions welcome — issues and pull requests in this repository.

## License

Apache-2.0 — see the [LICENSE](LICENSE) file (inherited from upstream; DaemonCores contributions under the same license).

---

<p>
  <strong align="left">Made with ⭐ by the DaemonCores community</strong>
  <a href="https://github.com/DaemonCores/traefik-proxmox-provider/wiki"><img align="right" src="https://img.shields.io/badge/Wiki-FFFFFF?style=for-the-badge&logoColor=white" alt="Documentation"/></a>
  <a href="https://github.com/orgs/DaemonCores/discussions"><img align="right" src="https://img.shields.io/badge/Community-000000?style=for-the-badge&logoColor=white" alt="Community"/></a>
  <a href="https://github.com/DaemonCores/traefik-proxmox-provider"><img align="right" src="https://img.shields.io/badge/Traefik_Proxmox_Provider-A81D33?style=for-the-badge&logo=traefikproxy&logoColor=white" alt="Traefik Proxmox Provider"/></a>
</p>