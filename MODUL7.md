# Modul 7 — RESTful API CRUD dengan Fiber v3 dan GORM

## Checklist persyaratan

| Persyaratan Modul 7 | Status |
|---|---|
| Koneksi MySQL melalui konfigurasi `.env` dan GORM | Implementasi ada; belum terverifikasi karena akun database memakai autentikasi GSSAPI |
| Model GORM mahasiswa serta relasi program studi | Implementasi selesai |
| AutoMigrate dan data awal program studi | Implementasi selesai; belum dapat dieksekusi pada database lokal |
| CRUD mahasiswa, pagination, pencarian, filter angkatan, dan preload relasi | Implementasi selesai |
| Pemisahan handler, service, dan repository | Implementasi selesai |
| CRUD matakuliah mengikuti lapisan handler/service/repository | Implementasi selesai |
| `GET /api/program-studi` dengan jumlah mahasiswa | Implementasi selesai; mahasiswa yang di-soft-delete tidak dihitung |
| Soft delete mahasiswa dan endpoint pemulihan | Implementasi selesai |
| Pencatatan kueri SQL GORM pada level Info | Implementasi selesai |
| Endpoint, route, header, contoh request dan response | Tercantum di bawah; route diuji tanpa koneksi DB |
| Perbandingan baris kode dan jawaban pertanyaan pembahasan | Tercantum di bawah |
| `go test -count=1 ./...`, `go vet ./...`, dan `go build ./...` | Lulus |
| Migrasi dan pengujian runtime dengan MariaDB | **Belum terverifikasi; Modul 7 belum dapat dinyatakan 100% selesai. Perubahan akun/database tidak disetujui.** |

### Pemeriksaan database

Tidak ada akun, password, kredensial, data, atau skema database yang diubah.

- `.env` berisi `DB_HOST=127.0.0.1`, `DB_PORT=3306`, `DB_NAME=pemweb2_fiber`, `DB_USER=root`, dan `DB_PASS` kosong.
- Session HeidiSQL bernama `Laragon.MySQL` juga menggunakan host `127.0.0.1`, port `3306`, user `root`, tanpa database awal yang dipilih, dan memiliki password tersimpan. Nilai password tidak dibaca atau ditampilkan.
- Greeting protokol pada port tersebut mengidentifikasi server sebagai MariaDB `12.3.3-MariaDB`. Database `pemweb2_fiber` ada pada data directory service MariaDB 12.3 dan tidak ada pada data directory Laragon MySQL 8.4.
- Service Windows `MariaDB` berjalan dari instalasi `C:\Program Files\MariaDB 12.3`. Instance MySQL Laragon 8.4 terpisah; port classic `3306` dipakai oleh service MariaDB, sedangkan proses Laragon saat pemeriksaan tidak memiliki listener classic MySQL pada port itu.
- Aplikasi Go gagal sebelum migrasi karena server meminta plugin client `auth_gssapi_client`, yang tidak didukung `go-sql-driver/mysql`. Klien MariaDB lokal berhasil memuat plugin GSSAPI setelah memakai path instalasi yang benar, tetapi login `root` tanpa password dari `.env` ditolak (`1045`); query hanya-baca tidak berhasil dijalankan.

Akun `root` lama tidak dapat dipakai langsung oleh driver Go dengan metode autentikasi saat ini. Solusi yang direkomendasikan adalah akun aplikasi baru dengan plugin autentikasi yang didukung driver serta hak minimum hanya pada `pemweb2_fiber`. Rencana pembuatan akun dan grant tersebut sudah dijelaskan kepada pengguna, tetapi pengguna memilih **jangan mengubah database**. Karena itu akun tidak dibuat, `.env` tidak diubah, AutoMigrate tidak dijalankan, dan modul belum boleh dinyatakan 100% selesai.

## Konfigurasi dan header

1. Pastikan database `pemweb2_fiber` sudah tersedia pada instance yang dipilih dan konfigurasi `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, serta `DB_PASS` di `.env` mengarah ke akun yang dapat diautentikasi oleh `go-sql-driver/mysql`. Jangan menganggap session HeidiSQL `Laragon.MySQL` menunjuk ke MySQL 8.4 Laragon hanya berdasarkan nama session.
2. `API_KEY` di `.env` wajib diisi. Semua endpoint API menggunakan header:

   ```text
   X-Api-Key: {{api_key}}
   Accept: application/json
   ```

3. Untuk request dengan body, tambahkan `Content-Type: application/json`. Jangan menyalin atau membagikan nilai rahasia `.env`.
4. Jalankan `go run ./cmd/api`. Jika koneksi berhasil, aplikasi menjalankan `AutoMigrate` untuk `ProgramStudi`, `Mahasiswa`, dan `Matakuliah`, lalu mengisi data awal program studi jika tabel masih kosong, sebelum membuka HTTP server.

## Route dan contoh request/response

Semua path di bawah relatif terhadap `base_url = http://localhost:3000/api`. Contoh response bersifat ilustratif; ID, jumlah data, dan timestamp bergantung pada database. Endpoint koleksi ditulis dengan trailing slash sesuai route yang didaftarkan.

### Autentikasi

Setiap endpoint pada tabel di bawah memerlukan `X-Api-Key`. Contoh uji:

```http
GET {{base_url}}/program-studi
Accept: application/json
```

Tanpa API key atau dengan nilai yang salah, response `401 Unauthorized`:

```json
{
  "sukses": false,
  "pesan": "API key tidak valid atau tidak dikirim"
}
```

### Mahasiswa

| Method dan path | Parameter | Request dan response sukses |
|---|---|---|
| `GET /mahasiswa/` | Query opsional: `halaman` (default 1), `per_halaman` (default 10, rentang efektif 1–100), `cari` (nama atau NIM), `angkatan` | Request: `GET {{base_url}}/mahasiswa/?halaman=1&per_halaman=5&cari=Eko&angkatan=2025`. Response `200`: `{"sukses":true,"data":[{"id":1,"program_studi_id":1,"program_studi":{"id":1,"kode":"TK","nama":"Teknik Komputer","jenjang":"S1","dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"},"nim":"H1A125001","nama":"Eko Saputra","email":"eko@example.com","angkatan":2025,"ipk":3.55,"aktif":true,"dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"}],"meta":{"halaman":1,"per_halaman":5,"total":1,"total_halaman":1}}` |
| `POST /mahasiswa/` | Body JSON | Request: `{"program_studi_id":1,"nim":"H1A125001","nama":"Eko Saputra","email":"eko@example.com","angkatan":2025,"ipk":3.55,"aktif":true}`. Response `201`: `{"sukses":true,"pesan":"Data mahasiswa berhasil dibuat","data":{"id":1,"program_studi_id":1,"nim":"H1A125001","nama":"Eko Saputra","email":"eko@example.com","angkatan":2025,"ipk":3.55,"aktif":true,"dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"}}` |
| `GET /mahasiswa/:id` | Path `id`: bilangan positif | Request: `GET {{base_url}}/mahasiswa/1`. Response `200`: `{"sukses":true,"data":{"id":1,"program_studi_id":1,"program_studi":{"id":1,"kode":"TK","nama":"Teknik Komputer","jenjang":"S1","dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"},"nim":"H1A125001","nama":"Eko Saputra","email":"eko@example.com","angkatan":2025,"ipk":3.55,"aktif":true,"dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"}}` |
| `PUT /mahasiswa/:id` | Path `id`; body JSON parsial | Request: `PUT {{base_url}}/mahasiswa/1`, body `{"ipk":3.85,"aktif":false}`. Response `200`: `{"sukses":true,"pesan":"Data mahasiswa berhasil diperbarui","data":{"id":1,"program_studi_id":1,"program_studi":{"id":1,"kode":"TK","nama":"Teknik Komputer","jenjang":"S1","dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"},"nim":"H1A125001","nama":"Eko Saputra","email":"eko@example.com","angkatan":2025,"ipk":3.85,"aktif":false,"dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"}}` |
| `DELETE /mahasiswa/:id` | Path `id` | Request: `DELETE {{base_url}}/mahasiswa/1`. Response `200`: `{"sukses":true,"pesan":"Data mahasiswa berhasil dihapus"}`. Record di-soft-delete, bukan dihapus permanen. |
| `POST /mahasiswa/:id/restore` | Path `id` harus menunjuk record yang di-soft-delete | Request: `POST {{base_url}}/mahasiswa/1/restore`. Response `200`: `{"sukses":true,"pesan":"Data mahasiswa berhasil dipulihkan","data":{"id":1,"program_studi_id":1,"program_studi":{"id":1,"kode":"TK","nama":"Teknik Komputer","jenjang":"S1","dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"},"nim":"H1A125001","nama":"Eko Saputra","email":"eko@example.com","angkatan":2025,"ipk":3.85,"aktif":false,"dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"}}` |

Saat memperbarui mahasiswa, handler menerima perubahan untuk `nama`, `email`, `angkatan`, `ipk`, `program_studi_id`, dan `aktif`. Field string kosong serta angka nol/negatif diabaikan; `aktif: false` tetap diterapkan. `nim` dipakai saat membuat mahasiswa, tetapi tidak diubah oleh handler update.

### Matakuliah

| Method dan path | Parameter | Request dan response sukses |
|---|---|---|
| `GET /matakuliah/` | Query opsional: `halaman` (default 1), `per_halaman` (default 10, rentang efektif 1–100), `cari` (kode atau nama), `semester` | Request: `GET {{base_url}}/matakuliah/?halaman=1&per_halaman=5&cari=Algoritma&semester=1`. Response `200`: `{"sukses":true,"data":[{"id":1,"kode":"IF101","nama":"Algoritma dan Pemrograman","sks":3,"semester":1,"dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"}],"meta":{"halaman":1,"per_halaman":5,"total":1,"total_halaman":1}}` |
| `POST /matakuliah/` | Body JSON; `kode` dan `nama` wajib; `sks` 1–6; `semester` 1–14; kode unik | Request: `POST {{base_url}}/matakuliah/`, body `{"kode":"IF101","nama":"Algoritma dan Pemrograman","sks":3,"semester":1}`. Response `201`: `{"sukses":true,"pesan":"Data matakuliah berhasil dibuat","data":{"id":1,"kode":"IF101","nama":"Algoritma dan Pemrograman","sks":3,"semester":1,"dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"}}` |
| `GET /matakuliah/:id` | Path `id`: bilangan positif | Request: `GET {{base_url}}/matakuliah/1`. Response `200`: `{"sukses":true,"data":{"id":1,"kode":"IF101","nama":"Algoritma dan Pemrograman","sks":3,"semester":1,"dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"}}` |
| `PUT /matakuliah/:id` | Path `id`; body JSON parsial berisi field yang akan diubah | Request: `PUT {{base_url}}/matakuliah/1`, body `{"sks":4}`. Response `200`: `{"sukses":true,"pesan":"Data matakuliah berhasil diperbarui","data":{"id":1,"kode":"IF101","nama":"Algoritma dan Pemrograman","sks":4,"semester":1,"dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"}}` |
| `DELETE /matakuliah/:id` | Path `id` | Request: `DELETE {{base_url}}/matakuliah/1`. Response `200`: `{"sukses":true,"pesan":"Data matakuliah berhasil dihapus"}`. Penghapusan matakuliah bersifat permanen. |

### Program studi

| Method dan path | Parameter | Request dan response sukses |
|---|---|---|
| `GET /program-studi` | Tidak ada | Request: `GET {{base_url}}/program-studi`. Response `200`: `{"sukses":true,"data":[{"id":1,"kode":"TK","nama":"Teknik Komputer","jenjang":"S1","jumlah_mahasiswa":1,"dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"},{"id":2,"kode":"IF","nama":"Informatika","jenjang":"S1","jumlah_mahasiswa":0,"dibuat_pada":"2026-11-05T15:00:00+07:00","diubah_pada":"2026-11-05T15:00:00+07:00"}]}` |

### Urutan pengujian Postman

1. Buat environment `base_url = http://localhost:3000/api` dan `api_key` dari konfigurasi lokal. Tambahkan header autentikasi pada seluruh request.
2. Pastikan akun database dapat diautentikasi oleh driver Go, lalu jalankan `go run ./cmd/api`. Lanjutkan hanya bila log menunjukkan koneksi dan migrasi berhasil.
3. Uji autentikasi dengan `GET {{base_url}}/program-studi`: tanpa API key diharapkan `401`; dengan API key lanjutkan pengujian.
4. Jalankan `GET /program-studi`, lalu `POST /mahasiswa/` menggunakan contoh body di atas. Gunakan NIM/email unik dan simpan ID dari response.
5. Uji daftar mahasiswa dengan `GET /mahasiswa/?halaman=1&per_halaman=5&cari=Eko&angkatan=2025`, detail dengan `GET /mahasiswa/{{mahasiswa_id}}`, dan pembaruan menggunakan contoh `PUT` di atas.
6. Jalankan `DELETE /mahasiswa/{{mahasiswa_id}}`, pastikan detail mengembalikan `404`, pulihkan dengan `POST /mahasiswa/{{mahasiswa_id}}/restore`, lalu pastikan detail kembali `200`.
7. Jalankan `POST /matakuliah/` menggunakan kode unik, simpan ID dari response, lalu uji list/filter, detail, update, dan delete dengan contoh di atas.

Kesalahan umum menggunakan bentuk `{"sukses":false,"pesan":"..."}`: `400` untuk body/ID tidak valid, `401` untuk API key hilang/salah, `404` untuk record tidak ditemukan, `422` untuk input/keunikan yang ditolak, dan `500` untuk kesalahan server/database.

## Pengujian yang masih memerlukan koneksi MariaDB

Persetujuan untuk perubahan akun/database belum diberikan. Sebelum pengujian di bawah dapat dilanjutkan, pemilik database perlu menyetujui solusi akun aplikasi terbatas atau menyiapkan sendiri akun yang kompatibel dan memperbarui konfigurasi lokal. Jangan mengganti password akun lama atau menjalankan AutoMigrate pada instance yang belum dipastikan.

- Startup aplikasi dengan akun yang berhasil diautentikasi dan eksekusi nyata `AutoMigrate` pada ketiga model.
- Verifikasi bahwa aplikasi diarahkan ke instance database yang memang dipilih (MariaDB 12.3 pada port 3306 saat ini, atau Laragon MySQL 8.4 yang merupakan instance berbeda).
- Verifikasi kolom, indeks unik, relasi, dan `deleted_at` yang terbentuk di database; pengulangan migrasi untuk memastikan aman/idempoten.
- Pembuatan data awal program studi dan verifikasi seed tidak menggandakan record pada startup berikutnya.
- Pengujian end-to-end setiap endpoint CRUD, pagination, pencarian/filter, serta respons ketika record tidak ada atau constraint unik dilanggar.
- Validasi `Preload("ProgramStudi")`, jumlah mahasiswa per program studi, dan pengecualian mahasiswa soft-delete dari hitungan.
- Pengujian persistensi soft delete/pemulihan setelah restart aplikasi.
- Pemeriksaan query SQL GORM yang tercetak saat migrasi dan endpoint berjalan.

## Perbandingan jumlah baris kode

Penghitungan untuk resource matakuliah pada project Laravel 13 dan Fiber ini memakai baris non-kosong/non-komentar dari komponen implementasi terkait:

| Implementasi | Bagian yang dihitung | Baris |
|---|---|---:|
| Fiber v3 | Handler 155, service 41, repository 42, model 11, DTO 6, route 6 | 261 |
| Laravel 13 | Controller 56, model 21, request validasi 46, resource 18, route 4 | 145 |

Implementasi Fiber pada project ini memakai sekitar 1,8 kali jumlah baris karena handler, service, dan repository serta pengolahan pagination/respons dibuat eksplisit. Laravel menyediakan abstraksi seperti `apiResource`, route model binding, Eloquent, Form Request, dan API Resource yang meringkas kode endpoint. Perbandingan ini hanya menghitung berkas terkait yang disebutkan, bukan kode framework.

## Jawaban pertanyaan pembahasan

1. Field `Aktif` pada DTO memakai pointer agar dapat membedakan nilai `false` yang dikirim dari field yang tidak dikirim. Saat pembuatan, field yang absen menggunakan default `true`.
2. `Save` menyimpan seluruh field model, termasuk nilai zero, dan dapat membuat record bila primary key belum ada. `Update` mengubah satu kolom; `Updates` mengubah beberapa kolom. `Updates` dengan struct biasanya melewati nilai zero, sedangkan map menyertakan nilai yang diberikan.
3. `AutoMigrate` tidak menghapus kolom yang tidak lagi dipakai dan tidak menggantikan migrasi berversi. Perubahan skema otomatis di produksi dapat sulit ditinjau/di-rollback dan berisiko mengunci atau mengubah tabel saat aplikasi berjalan.
4. `Preload` pada GORM dan `with()` pada Eloquent sama-sama eager loading relasi untuk menghindari N+1, umumnya dengan query relasi terpisah. GORM juga menyediakan join preload melalui `Joins`; Eloquent mendukung constraint dan variasi eager loading pada `with()`.
