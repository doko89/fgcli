# fgcli — FortiGate CLI for humans & AI agents

## Contract (AI-friendly)
- Semua stdout = **satu envelope JSON**: `{"ok":true,"data":...}` atau `{"ok":false,"error":{"code":"...","message":"..."}}`
- Log manusia ke **stderr**, exit code: `0` ok, `1` usage/input/validasi, `2` API/connect.
- Non-interaktif: destruktif wajib `--yes`, ada `--dry-run`, input JSON via `--from-stdin` / `--from-file`.

## Config
API key:
```bash
export FG_HOST="https://192.0.2.1" FG_API_KEY="<api-key>" FG_VDOM="root" FG_INSECURE=1
./fgcli doctor   # tes koneksi + status
```

## Profile (satu file: `~/.fgcli/config.yaml`, override via `FGCLI_CONFIG`)
```bash
export FG_API_KEY_FG1='...'   # secret via env, tidak tertulis di disk
./fgcli profile add fg1 --host https://172.28.29.1:15443 --api-key-env FG_API_KEY_FG1 --insecure
./fgcli profile use fg1    # switch profil aktif (alias: profile set fg1)
./fgcli profile list       # secret tidak pernah ditampilkan
./fgcli profile validate   # offline: placeholder key, host duplikat, insecure
./fgcli doctor --profile fg1
./fgcli --profile fg1 address list
```
Hanya `~/.fgcli/config.yaml` yang dibaca; file lama (`~/fgcli/config.yaml`,
`~/tmp/FG-CLI/config.ini`) diabaikan (ada hint di stderr bila masih ada).
`api_key` placeholder (`will_be_available`, `changeme`, `xxx`, ...) langsung
ditolak sebelum request; `insecure:true` selalu memicu warning di stderr.
Precedence: flag (`--host/--api-key`) > env (`FG_*`) > profile aktif.
File ditulis chmod `0600`. Kompatibel baca: `FG_API_KEY`/`--token`/`token:` lama tetap dibaca sebagai API key.

## NAT (P0)
```bash
./fgcli vip list --filter 103.165      # semua field ikut dicari (extip, mappedip range)
./fgcli vip get OXYGEN_172.23.9.2_20
./fgcli ippool list                    # outbound SNAT pool
./fgcli snat list                      # central-snat-map ([] = pakai poolname di policy)
./fgcli policy list --filter VIP_MAIL  # policy yang refer ke VIP (dstaddr)
./fgcli address list --contains 192.168.1.43
```
`--filter`/`--contains` (alias) substring case-insensitive di semua `list`.
`vip/ippool` full CRUD: `create/update` via flag utama atau `--from-stdin` JSON,
`delete` wajib `--yes`, semua dukung `--dry-run`.

## Service & interface (P1)
```bash
./fgcli service get HTTPS              # -> tcp-portrange 443
./fgcli service list --filter 822
./fgcli service-group get "Email Access"
./fgcli interface list                 # slim: name, ip, status, type, alias
./fgcli interface get wan1
./fgcli zone list
```
`service`/`service-group` full CRUD; `interface`/`zone` read-only
(platform config, tidak diprovision via tool ini).

## Monitor Batch 1
```bash
./fgcli system license        # forticare status/account/support + fortiguard connected
./fgcli system fortiguard     # service-communication-stats per layanan
./fgcli vpn ipsec             # 3 tunnel (RT_SitetoIDXREP, HO_TO_DC_TUNNEL, TO_ITCH) + bytes
./fgcli routing status        # total_lines 28 (ipv4 28, ipv6 0)
./fgcli system ntp            # 4 server fortiguard, reachable
./fgcli system dns            # [] di box ini (acquired-dns kosong)
./fgcli system dhcp           # [] — HTTP 424 = tidak ada interface DHCP, bukan error
./fgcli system snmp           # status disable
```

## User & Security Batch 2
```bash
./fgcli user firewall          # 2 login (FZ_14A, ahmad) + grup + durasi
./fgcli user banned            # [] di box ini
./fgcli security ips           # 18 anomali (tcp_syn_flood, ...)
./fgcli security waf list      # 11 kelas; get 30000000 → SQL Injection
./fgcli security dlp get 1     # builtin-patterns, 18 entries
./fgcli security proxy-pac --out proxy.pac   # 424 → file kosong (tidak dikonfigurasi)
./fgcli service voip get default             # voipd, sip enable
./fgcli network dnsfilter list               # [] di box ini
```

## Switch/WiFi/FortiView/Log Batch 3+4
```bash
./fgcli switch status      # [] — tidak ada managed switch
./fgcli wifi ap            # semua counter 0 (nol = data, tidak di-drop)
./fgcli fortiview          # summary start/end, details kosong
./fgcli log memory dns     # [] — body kosong = tidak ada log
./fgcli log disk ips       # tipe bebas, server yang validasi (404 = tipe salah)
```

## Diagnose (P2)
```bash
./fgcli raw get /cmdb/firewall/vip                 # path boleh tanpa prefix /api/v2
./fgcli raw get /monitor/system/status
./fgcli raw post /cmdb/firewall/address --data '{"name":"x","subnet":"10.0.0.9/32"}' --dry-run
./fgcli trace --dst 103.118.129.220 --dport 443    # vip match + policy terurut
./fgcli trace --src 8.8.8.8 --dst 103.165.33.70 --dport 443  # negatif: bukan DNAT box ini
```
`raw`: escape hatch terautentikasi, `--vdom` hanya dikirim bila eksplisit
(monitor path tidak butuh vdom). `trace`: korelasi vip (extip+dport) +
policy (src/dst/service, referensi nama maupun subnet `IP MASK`/CIDR);
policy spesifik diurut dulu, `disable` ditandai di `via`.

## Help, grup, policy lifecycle, backup
```bash
./fgcli --help                # == fgcli help == fgcli -h (tanpa auth, exit 0)
./fgcli vip --help            # tanpa auth, exit 0 (juga: fgcli help vip)
./fgcli vipgrp list           # [] di box ini (endpoint valid, kosong)
./fgcli addrgrp list          # 6 grup (BO_RT, G Suite, ...)
./fgcli policy create --policyid 99 --name TRIAL --srcintf internal4 \
  --dstintf wan1 --srcaddr all --dstaddr all --service HTTPS --dry-run
./fgcli policy enable|disable <id> | move <id> --before X | clone <s> --new-id N
./fgcli backup download --out fg.conf   # butuh user RW (readonly → 403)
./fgcli raw get /cmdb/system/interface/wan1   # detail interface penuh
```
Catatan: `ippoolgrp` 400 di firmware 7.4.9 (dengan/tanpa vdom) → verb tidak
dibuat; pakai `raw` bila firmware lain support. Write API (create/update/
move/delete) terverifikasi via `--dry-run` + httptest karena user readonly.

## Pakai
```bash
go build -o fgcli ./cmd/fgcli
make install   # -> ~/.local/bin/fgcli (pastikan di PATH: export PATH="$HOME/.local/bin:$PATH")
./fgcli version
./fgcli --pretty doctor   # --pretty = JSON indent untuk manusia; default tetap compact
./fgcli vpn ssl get       # alias: raw get cmdb/vpn.ssl/settings
./fgcli user group add-member VPN doko_baru --dry-run
./fgcli user local create doko_baru --password xxx --group VPN --dry-run
./fgcli doctor
./fgcli system status
./fgcli address list --vdom root
./fgcli address get lan-net
./fgcli address create --name lan-net --subnet 10.0.0.0/24 --comment "via ai"
echo '{"name":"lan-net","subnet":"10.0.0.0/24"}' | ./fgcli address create --from-stdin --dry-run
./fgcli address delete lan-net --yes
./fgcli policy list
./fgcli policy delete 5 --yes
```

## Struktur DDD
```
cmd/fgcli              composition root (wiring)
internal/domain        entity + Repository interface (address, policy, svc)
internal/application   usecase (validasi + orkestrasi)
internal/infrastructure forte: fortios client (timeout 15s, retry 3x), config, output JSON
internal/interfaces/cli dispatcher CLI (stdout JSON, stderr log)
```

## Tambah resource baru
1. `internal/domain/<x>/` entity + `Repository` interface.
2. `internal/application/<x>/` usecase.
3. `internal/infrastructure/fortios/` implementasi repo (`/api/v2/cmdb/...`).
4. `internal/interfaces/cli/` tambah verb + wiring di `cmd/fgcli/main.go`.
