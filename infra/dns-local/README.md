# Local DNS Setup untuk *.dejavu.local

Panduan setup wildcard subdomain untuk development lokal.

## macOS

### 1. Install dnsmasq

```bash
brew install dnsmasq
```

### 2. Configure dnsmasq

```bash
# Copy config
sudo cp dnsmasq.conf /usr/local/etc/dnsmasq.conf

# Atau tambahkan manual
echo "address=/.dejavu.local/127.0.0.1" | sudo tee -a /usr/local/etc/dnsmasq.conf
```

### 3. Start dnsmasq

```bash
sudo brew services start dnsmasq
```

### 4. Configure resolver

```bash
sudo mkdir -p /etc/resolver
sudo tee /etc/resolver/dejavu.local <<EOF
nameserver 127.0.0.1
EOF
```

### 5. Test

```bash
ping test.dejavu.local
# Should resolve to 127.0.0.1
```

## Linux (Ubuntu/Debian)

### 1. Install dnsmasq

```bash
sudo apt update
sudo apt install dnsmasq
```

### 2. Configure dnsmasq

```bash
sudo cp dnsmasq.conf /etc/dnsmasq.conf

# Atau tambahkan manual
echo "address=/.dejavu.local/127.0.0.1" | sudo tee -a /etc/dnsmasq.conf
```

### 3. Restart dnsmasq

```bash
sudo systemctl restart dnsmasq
sudo systemctl enable dnsmasq
```

### 4. Configure NetworkManager (jika menggunakan)

```bash
sudo tee /etc/NetworkManager/conf.d/dnsmasq.conf <<EOF
[main]
dns=dnsmasq
EOF

sudo systemctl restart NetworkManager
```

### 5. Test

```bash
dig test.dejavu.local
# Should return 127.0.0.1
```

## Windows

### Metode 1: Acrylic DNS Proxy

1. Download Acrylic DNS Proxy dari https://mayakron.altervista.org/support/acrylic/Home.htm
2. Install dan jalankan sebagai service
3. Edit `AcrylicHosts.txt`:
   ```
   127.0.0.1 *.dejavu.local
   ```
4. Restart Acrylic service
5. Set DNS di Network Settings ke `127.0.0.1`

### Metode 2: Edit hosts file (manual per subdomain)

1. Buka Notepad as Administrator
2. Edit `C:\Windows\System32\drivers\etc\hosts`
3. Tambahkan:
   ```
   127.0.0.1 dejavu.local
   127.0.0.1 app1.dejavu.local
   127.0.0.1 app2.dejavu.local
   ```
4. Save dan restart browser

## Docker Desktop (All Platforms)

Jika menggunakan Docker Desktop, tambahkan di `docker-compose.yml`:

```yaml
services:
  traefik:
    ports:
      - "80:80"
    labels:
      - "traefik.http.routers.wildcard.rule=HostRegexp(`{subdomain:.+}.dejavu.local`)"
```

## Verification

Test dengan salah satu cara:

```bash
# Ping
ping anything.dejavu.local

# nslookup
nslookup random123.dejavu.local

# curl
curl http://test.dejavu.local
```

Semua harus resolve ke `127.0.0.1`.

## Troubleshooting

### macOS: dnsmasq not working

```bash
# Check if running
sudo brew services list

# Check logs
tail -f /usr/local/var/log/dnsmasq.log

# Restart
sudo brew services restart dnsmasq
```

### Linux: dnsmasq conflicts with systemd-resolved

```bash
# Stop systemd-resolved
sudo systemctl stop systemd-resolved
sudo systemctl disable systemd-resolved

# Update /etc/resolv.conf
echo "nameserver 127.0.0.1" | sudo tee /etc/resolv.conf
```

### Windows: DNS Cache

```powershell
# Clear DNS cache
ipconfig /flushdns
```

