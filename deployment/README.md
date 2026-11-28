# Persiapan deployment Praktikum Modul 10

Konfigurasi ini menargetkan Ubuntu Server VMware. Ganti `<IP_SERVER>` dan `<DOMAIN>` dengan nilai aktual; tidak ada alamat server yang diasumsikan. Jangan menyalin `.env` aktif atau memasukkan rahasia ke Git.

## 1. Build artefak di Windows

Dari root project:

```powershell
go test ./...
go vet ./...
$env:CGO_ENABLED = "0"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
New-Item -ItemType Directory -Force bin | Out-Null
go build -ldflags="-s -w" -o bin/api ./cmd/api
Remove-Item Env:CGO_ENABLED, Env:GOOS, Env:GOARCH
```

Binary `bin/api` adalah Linux AMD64 dan tidak dapat dijalankan di Windows. Frontend dibangun terpisah dari `frontend-pemweb2` dengan `npm run build`; file statis ada di `frontend-pemweb2/dist/`.

## 2. Siapkan server dan unggah file

Pastikan Ubuntu Server dapat mengakses database yang sudah disiapkan. Sebelum migrasi pertama, buat backup database dan tinjau file migrasi. Jangan menjalankan `migrate:fresh`, `db:wipe`, atau reset database.

Salin binary, source Laravel, dan hasil build frontend dari komputer Windows ke server. Contoh transfer binary (jalankan di Windows, ganti alamat dan nama pengguna):

```powershell
ssh <USER>@<IP_SERVER> "mkdir -p /tmp/pemweb2-laravel"
scp .\bin\api <USER>@<IP_SERVER>:/tmp/pemweb2-api
scp -r .\latihan-laravel\app .\latihan-laravel\bootstrap .\latihan-laravel\config .\latihan-laravel\database\migrations .\latihan-laravel\database\seeders .\latihan-laravel\public .\latihan-laravel\resources .\latihan-laravel\routes .\latihan-laravel\artisan .\latihan-laravel\composer.json .\latihan-laravel\composer.lock .\latihan-laravel\.env.example <USER>@<IP_SERVER>:/tmp/pemweb2-laravel/
scp -r .\frontend-pemweb2\dist <USER>@<IP_SERVER>:/tmp/pemweb2-dist
scp .\deployment\pemweb2-api.service .\deployment\pemweb2-laravel.service .\deployment\nginx-pemweb2.conf .\deployment\fiber-api.env.example <USER>@<IP_SERVER>:/tmp/
```

Perintah di atas sengaja tidak menyalin `.env`, `vendor`, atau file database lokal.

Di Ubuntu Server, instal paket sesuai versi Ubuntu/PHP dan driver database yang dipakai:

```bash
sudo apt update
sudo apt install -y nginx php-cli php-mbstring php-xml php-curl php-sqlite3 php-mysql unzip composer
sudo useradd --system --home /opt/pemweb2 --shell /usr/sbin/nologin pemweb2
sudo install -d -o pemweb2 -g pemweb2 /opt/pemweb2/fiber-api
sudo install -d -o www-data -g www-data /opt/pemweb2/laravel
sudo install -d -o www-data -g www-data /opt/pemweb2/laravel/storage/framework/cache/data /opt/pemweb2/laravel/storage/framework/sessions /opt/pemweb2/laravel/storage/framework/views /opt/pemweb2/laravel/storage/logs /opt/pemweb2/laravel/bootstrap/cache
sudo install -d -o www-data -g www-data /var/www/pemweb2/frontend
sudo install -o root -g root -m 0755 /tmp/pemweb2-api /opt/pemweb2/fiber-api/api
sudo cp -a /tmp/pemweb2-laravel/. /opt/pemweb2/laravel/
sudo chown -R www-data:www-data /opt/pemweb2/laravel
sudo cp -a /tmp/pemweb2-dist/. /var/www/pemweb2/frontend/dist/
sudo chown -R www-data:www-data /opt/pemweb2/laravel/storage /opt/pemweb2/laravel/bootstrap/cache
```

Jika Laravel memakai MySQL, instal ekstensi `php-mysql` dan konfigurasi database di berkas environment. Jika Fiber memakai MySQL, pastikan server DB/firewall mengizinkan koneksi hanya dari host API.

## 3. Konfigurasi Laravel

```bash
cd /opt/pemweb2/laravel
sudo -u www-data cp .env.example .env
sudo -u www-data composer install --no-dev --optimize-autoloader --no-interaction
sudo -u www-data php artisan key:generate --force
sudo -u www-data nano .env
sudo chmod 0640 .env
```

Atur `APP_ENV=production`, `APP_DEBUG=false`, `APP_URL=https://<DOMAIN>`, `FRONTEND_URL=https://<DOMAIN>` (atau origin frontend yang tepat), `DB_CONNECTION=mysql`, `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD`, dan konfigurasi cache/session yang sesuai. Siapkan database dan akun DB dengan izin minimum terlebih dahulu; gunakan ekstensi PHP database yang sesuai. Simpan `.env` hanya pada server dengan izin baca terbatas:

Setelah backup dan review skema:

```bash
sudo -u www-data php artisan migrate --force
sudo -u www-data php artisan config:cache
sudo -u www-data php artisan route:cache
sudo -u www-data php artisan view:cache
```

Service Artisan port 8000 disediakan untuk kebutuhan praktikum reverse proxy. Untuk produksi sesungguhnya, gunakan PHP-FPM dan konfigurasi Nginx Laravel standar, bukan `artisan serve`.

## 4. Konfigurasi Fiber

Buat file environment di server, bukan di repository:

```bash
sudo install -d -m 0750 -o root -g pemweb2 /etc/pemweb2
sudo install -m 0640 -o root -g pemweb2 /tmp/fiber-api.env.example /etc/pemweb2/fiber-api.env
sudo nano /etc/pemweb2/fiber-api.env
```

Ganti semua nilai `replace-with-*` dengan nilai rahasia yang dibuat di server (contoh `openssl rand -hex 32` untuk API key dan JWT secret), dan sesuaikan `CORS_ALLOWED_ORIGINS=https://<DOMAIN>` (origin frontend spesifik, bukan `*`) serta koneksi DB. Jika frontend dan API memakai host yang sama melalui Nginx, browser tidak memerlukan CORS untuk request same-origin.

```bash
sudo install -m 0644 /tmp/pemweb2-api.service /etc/systemd/system/pemweb2-api.service
sudo systemctl daemon-reload
sudo systemctl enable --now pemweb2-api
sudo systemctl status pemweb2-api --no-pager
sudo journalctl -u pemweb2-api -n 100 --no-pager
```

Fiber menjalankan `AutoMigrate` dan data awal saat startup. Periksa skema/backup terlebih dahulu; jangan restart service terhadap DB produksi yang tidak disiapkan atau tidak dicadangkan.

## 5. Nginx dan service Laravel

```bash
sudo install -m 0644 /tmp/pemweb2-laravel.service /etc/systemd/system/pemweb2-laravel.service
sudo systemctl daemon-reload
sudo systemctl enable --now pemweb2-laravel
sudo install -m 0644 /tmp/nginx-pemweb2.conf /etc/nginx/sites-available/pemweb2
sudo ln -s /etc/nginx/sites-available/pemweb2 /etc/nginx/sites-enabled/pemweb2
sudo nginx -t
sudo systemctl reload nginx
```

`proxy_pass` dengan trailing slash menghapus prefix `/laravel/` atau `/fiber/` sebelum meneruskan request; endpoint aplikasi tetap menggunakan `/api/...`. Jika terdapat konfigurasi default Nginx yang bentrok, nonaktifkan hanya setelah memeriksa konfigurasi aktif. Jangan membuka port 8000/3000 ke jaringan publik; akseskan API melalui Nginx saja.

Untuk HTTPS pada domain yang benar-benar diarahkan ke VM, ubah `server_name _;` menjadi `server_name <DOMAIN>;`, jalankan `sudo nginx -t && sudo systemctl reload nginx`, pasang Certbot sesuai instruksi Ubuntu, lalu jalankan `sudo certbot --nginx -d <DOMAIN>`. Di jaringan VM lokal tanpa domain, uji HTTP hanya di jaringan terisolasi dan jangan anggap itu deployment publik yang aman.

## 6. Verifikasi dari Ubuntu dan Kali

Di Ubuntu:

```bash
curl -i http://127.0.0.1/laravel/api/status
curl -i http://127.0.0.1/fiber/api/status
sudo systemctl is-active pemweb2-api pemweb2-laravel nginx
sudo nginx -t
```

Dari Kali Linux, ganti `<IP_SERVER>`:

```bash
curl -i http://<IP_SERVER>/laravel/api/status
curl -i http://<IP_SERVER>/fiber/api/status
curl -i http://<IP_SERVER>/laravel/api/mahasiswa
curl -i http://<IP_SERVER>/fiber/api/mahasiswa
```

Status dan endpoint mahasiswa yang memerlukan token dapat diverifikasi dengan HTTP 200 dan 401. Setelah HTTPS aktif, ulangi dengan `https://<DOMAIN>` dan pastikan sertifikat valid. Jangan mengirim kredensial, API key, atau token melalui HTTP pada jaringan yang tidak tepercaya.

## 7. Backup database

Tentukan jadwal, retensi, enkripsi, lokasi salinan terpisah, dan uji pemulihan. Untuk MySQL, gunakan akun backup dengan izin minimum dan simpan dump di lokasi privat, bukan repository:

```bash
mysqldump --single-transaction -u <DB_BACKUP_USER> -p <DB_NAME> > <PATH_BACKUP>/<DB_NAME>-$(date +%F).sql
```

Verifikasi ukuran/exit status dump, salin backup terenkripsi ke media terpisah, dan uji restore berkala pada database non-produksi.
