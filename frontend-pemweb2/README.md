# Praktikum Modul 9 — Integrasi Backend dan Frontend

Frontend React + Vite untuk mengakses API Laravel atau Go Fiber. Pemilihan backend dan penyesuaian kontrak dilakukan melalui konfigurasi dan adapter frontend; backend tetap menggunakan kontraknya masing-masing.

## Instalasi dan menjalankan

Pastikan Node.js/npm tersedia. Jalankan setiap server pada terminal terpisah.

### Laravel

```powershell
Set-Location 'D:\kuliah\Pemweb2\latihan-laravel'
php artisan serve --host 127.0.0.1 --port 8000
```

Laravel memerlukan konfigurasi database dan data pengguna yang sudah tersedia.

### Fiber

```powershell
Set-Location 'D:\kuliah\Pemweb2'
go run ./cmd/api
```

Gunakan modul Fiber pada root karena pengujian, build, dan GitHub Actions Modul 10 mengacu pada `go.mod` root. Fiber memerlukan konfigurasi database pada `.env`; proses startup menjalankan migrasi dan data awal.

### Frontend

```powershell
Set-Location 'D:\kuliah\Pemweb2\frontend-pemweb2'
npm install
npm run dev
```

Buka URL yang ditampilkan Vite, biasanya `http://localhost:5173`.

## Mengganti backend

Salin `.env.example` menjadi `.env` untuk pengembangan lokal, lalu pilih backend:

```dotenv
VITE_API_BACKEND=laravel
```

Gunakan `laravel` atau `fiber`, lalu mulai ulang server Vite. URL kedua backend disimpan masing-masing pada `VITE_API_LARAVEL` dan `VITE_API_FIBER`. `src/api/kontrak.js` menyatukan perbedaan request, response, pagination, dan galat. Build produksi memakai URL relatif `/laravel/api` dan `/fiber/api` dari `.env.production`, sehingga dapat dilayani bersama melalui Nginx tanpa menetapkan IP server.

## Ringkasan endpoint yang digunakan

| Operasi | Laravel | Fiber |
| --- | --- | --- |
| Login | `POST /api/auth/login` | `POST /api/auth/login` |
| Profil | `GET /api/auth/profil` | `GET /api/auth/profil` |
| Logout | `POST /api/auth/logout` | `POST /api/auth/logout` |
| Daftar mahasiswa | `GET /api/mahasiswa` | `GET /api/mahasiswa` |
| Detail mahasiswa | `GET /api/mahasiswa/{id}` | `GET /api/mahasiswa/{id}` |
| Tambah mahasiswa | `POST /api/mahasiswa` | `POST /api/mahasiswa` |
| Hapus mahasiswa | `DELETE /api/mahasiswa/{id}` | `DELETE /api/mahasiswa/{id}` |

Endpoint mahasiswa memerlukan token Bearer. Izin tulisnya berbeda: Laravel memeriksa ability `mahasiswa:tulis`; Fiber mengizinkan admin/dosen menambah atau mengubah dan hanya admin menghapus. Endpoint matakuliah Fiber juga memakai `X-Api-Key`; operasi tulisnya memerlukan API key dan token Bearer dengan peran yang sesuai. Daftar program studi Fiber memakai API key saja. Laravel membatasi endpoint tulis matakuliah dengan ability yang sama.

## Perbandingan kontrak API

| Bagian | Laravel | Fiber | Penyesuaian frontend |
| --- | --- | --- | --- |
| Password login | `password` | `kata_sandi` | `buatPayloadLogin()` memakai field sesuai backend aktif |
| Pengguna pada login | `data.user`, nama field `name` | `data.pengguna`, nama field `nama` | Dinormalisasi menjadi `user.nama` |
| Profil | `data.user` | `data` berisi pengguna | Dinormalisasi menjadi objek profil yang sama |
| Token | Sanctum personal access token opaque | JWT HS256; TTL default 120 menit | Keduanya dikirim sebagai `Authorization: Bearer ...` |
| Nomor halaman request | `page` | `halaman` | `buatParameterDaftarMahasiswa()` memilih nama parameter |
| Pagination response | `data`, `links`, `meta.current_page/per_page/last_page/total` | `data`, `meta.halaman/per_halaman/total_halaman/total` | Dinormalisasi ke `pagination.halaman/perHalaman/totalHalaman/total` |
| Galat validasi | HTTP 422 dengan `errors[field]: [...]` | HTTP 422 dengan `galat: [{field,message}]` | Dinormalisasi ke `fieldErrors` |
| Login gagal | HTTP 422 dengan `errors` | HTTP 401 dengan `pesan` | Pesan dan galat field ditampilkan seragam |
| Validasi NIM | Panjang/unik; tidak mewajibkan pola prefix | Harus `H` diikuti angka, 5–20 karakter | Galat backend ditampilkan apa adanya |
| Program studi | `program_studi_id` harus ada pada database | `program_studi_id >= 1`; database tetap memeriksa relasi | Form menerima ID program studi yang sudah terdaftar |
| Bentuk detail mahasiswa | Resource, termasuk `program_studi` terpilih | Model GORM, termasuk relasi dan timestamp | Data inti digunakan bersama; field tambahan backend dipertahankan |

## Usulan kontrak API tunggal

Kontrak berikut dapat dipakai bila backend kelak diseragamkan. Implementasi frontend sekarang tidak mengubah kontrak backend.

### Login dan profil

Request login menggunakan `email` dan `password`. Respons login dan profil menggunakan wrapper yang konsisten:

```json
{
  "data": {
    "user": {
      "id": 1,
      "nama": "Nama pengguna",
      "email": "pengguna@example.com",
      "peran": "admin"
    },
    "token": "token",
    "token_type": "Bearer"
  }
}
```

Profil dapat menggunakan `data.user` tanpa `token`.

### Daftar dan pagination

Gunakan parameter `page`, `per_page`, `search`, dan `angkatan`. Bentuk respons:

```json
{
  "data": [],
  "pagination": {
    "current_page": 1,
    "per_page": 10,
    "total": 0,
    "last_page": 1
  }
}
```

### Galat

Gunakan status HTTP yang konsisten dan format:

```json
{
  "message": "Data yang dikirim tidak valid",
  "errors": {
    "nim": ["Format NIM tidak valid"]
  }
}
```

Dengan kontrak tersebut, klien tidak perlu menebak bentuk respons, nama field password, atau struktur pagination.

## Pengujian

- `npm run lint` dan `npm run build` lulus; `go test ./internal/middleware ./cmd/api` juga lulus.
- `php artisan migrate:status --no-ansi` berhasil tanpa mengubah database; seluruh migrasi Laravel berstatus `Ran`.
- Laravel `GET /api/status` menghasilkan HTTP 200. Preflight `OPTIONS /api/mahasiswa` dari `http://localhost:5173` menghasilkan HTTP 204 dengan metode dan header yang diminta.
- Browser mengirim permintaan mahasiswa dengan Bearer token uji yang tidak valid; Laravel merespons HTTP 401 dalam 62 ms dengan `Access-Control-Allow-Origin: http://localhost:5173`. Klien menghapus token dan mengarahkan kembali ke `/login`.
- Uji browser dari origin yang tidak diizinkan (`http://127.0.0.1:5173`) gagal sesuai harapan: preflight ditolak browser karena origin respons tidak cocok. Konfigurasi CORS tidak diubah selama uji.
- Uji tambah, validasi HTTP 422, konfirmasi hapus, pembatalan, dan refresh daftar menggunakan fixture HTTP lokal; fixture tidak menulis ke database backend.

Login berhasil serta daftar/tambah/hapus terhadap backend dengan token valid belum diverifikasi karena tidak ada kredensial uji yang diberikan. Server Fiber tidak dijalankan: startup-nya menjalankan `AutoMigrate` dan dapat membuat data awal; pengujian Fiber yang dijalankan terbatas pada test Go agar database tidak berubah tanpa verifikasi.
