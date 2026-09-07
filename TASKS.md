# TASKS — Gap fgcli vs example/FG-CLI

Hasil ukur: `example/FG-CLI` (Python, menu interaktif, token-only) punya
~25 endpoint read-only yang belum ada verb native di fgcli v0.2.0.
Semua di bawah ini read-only (GET) kecuali backup (sudah ada).
Usulan command mengikuti konvensi fgcli: non-interaktif, envelope JSON,
`--filter` gratis via matcher generik.

Catatan bug referensi (jangan dicopy): key `dns-filter` vs `dns_filter`
salah rujuk (menu Network→DNS Filter selalu error); key `config_backup`
tidak dipakai (backup asli lewat `backup2.py`); URL di `backup2.py`
masih placeholder.

## Batch 1 — Monitor inti (P1)
- [x] `system license` — GET monitor/license/status
- [x] `system fortiguard` — GET monitor/fortiguard/service-communication-stats
- [x] `vpn ipsec` — GET monitor/vpn/ipsec (audit tunnel DEV/HO butuh ini)
- [x] `routing status` — GET monitor/router/statistics
- [x] `system ntp` — GET monitor/system/ntp/status
- [x] `system dns` — GET monitor/system/acquired-dns
- [x] `system dhcp` — GET monitor/system/interface/dhcp-status
- [x] `system snmp` — GET cmdb/system.snmp/sysinfo

## Batch 2 — User & Security (P1)
- [x] `user firewall` — GET monitor/user/firewall
- [x] `user banned` — GET monitor/user/banned
- [x] `security ips` — GET monitor/ips/anomaly
- [x] `security waf get <id>` — GET cmdb/waf/main-class/{id}
- [x] `security dlp get <id>` — GET cmdb/dlp/filepattern/{id}
- [x] `security proxy-pac` — GET monitor/webproxy/pacfile/download (simpan file, pola backup)
- [x] `service voip get <name>` — GET cmdb/voip/profile/{name}
- [x] `network dnsfilter get <id>` — GET cmdb/dnsfilter/domain-filter/{id}

## Batch 3 — Switch/WiFi & FortiView (P2)
- [x] `switch status` — GET monitor/switch-controller/managed-switch/status
- [x] `wifi ap` — GET monitor/wifi/ap_status
- [x] `fortiview` — GET monitor/fortiview/statistics (respons besar, pastikan limit 8MB cukup)

## Batch 4 — Logging (P2)
- [x] `log fortianalyzer <type>` — GET log/fortianalyzer/{type}/raw
- [x] `log forticloud <type>` — GET log/forticloud/{type}/raw
- [x] `log memory <type>` — GET log/memory/{type}/raw
- [x] `log disk <type>` — GET log/disk/{type}/raw
- [x] Tipe log (argumen bebas, divalidasi server): virus, webfilter, waf, ips, anomaly, app-ctrl, emailfilter, dlp, voip, gtp, dns, ssh, ssl, cifs, file-filter (argumen, bukan submenu)

## Sudah covered (jangan duplikat)
- Backup config → `backup download` (referensi: backup2.py)
- system status → `system status` / `doctor`
- CMDB firewall/service/system CRUD → fgcli lebih lengkap dari referensi
  (referensi tidak punya write apa pun)

## Keunggulan fgcli yg tidak ada di referensi (pertahankan)
- JSON envelope `{ok,data,error}`, `--dry-run`, `--filter`, `--from-stdin`
- Auth API key + multi-profile
- `trace` korelasi NAT/policy, `raw` escape hatch, help tanpa auth
