const backend = {
  laravel: {
    baseURL: import.meta.env.VITE_API_LARAVEL,
    passwordField: 'password',
    pageField: 'page',
  },
  fiber: {
    baseURL: import.meta.env.VITE_API_FIBER,
    passwordField: 'kata_sandi',
    pageField: 'halaman',
  },
}

export const backendAktif = import.meta.env.VITE_API_BACKEND
export const konfigurasiApi = backend[backendAktif]

if (!konfigurasiApi?.baseURL) {
  throw new Error(`Backend API "${backendAktif}" tidak dikenal atau URL belum dikonfigurasi`)
}

export function buatPayloadLogin(email, kataSandi) {
  return {
    email,
    [konfigurasiApi.passwordField]: kataSandi,
  }
}

export function buatParameterDaftarMahasiswa({
  halaman,
  perHalaman,
  cari,
  angkatan,
}) {
  return {
    [konfigurasiApi.pageField]: halaman,
    per_halaman: perHalaman,
    cari,
    angkatan,
  }
}

function objekWajib(nilai, nama) {
  if (!nilai || typeof nilai !== 'object' || Array.isArray(nilai)) {
    throw new Error(`Respons API tidak memiliki objek ${nama} yang valid`)
  }

  return nilai
}

function normalisasiPengguna(nilai, nama) {
  const pengguna = objekWajib(nilai, nama)
  const namaPengguna = pengguna.nama ?? pengguna.name

  if (pengguna.id === undefined || !namaPengguna || !pengguna.email) {
    throw new Error(`Respons API ${nama} tidak memuat identitas pengguna yang lengkap`)
  }

  return {
    id: pengguna.id,
    nama: namaPengguna,
    email: pengguna.email,
    peran: pengguna.peran,
  }
}

export function normalisasiLogin(respons) {
  const data = objekWajib(objekWajib(respons, 'login').data, 'login.data')
  const token = data.token

  if (typeof token !== 'string' || token.length === 0) {
    throw new Error('Respons API login tidak memuat token')
  }

  return {
    user: normalisasiPengguna(data.user ?? data.pengguna, 'login.user'),
    token,
    tokenType: data.tipe ?? 'Bearer',
  }
}

export function normalisasiProfil(respons) {
  const data = objekWajib(objekWajib(respons, 'profil').data, 'profil.data')
  return normalisasiPengguna(data.user ?? data, 'profil.user')
}

export function normalisasiDaftarMahasiswa(respons) {
  const data = objekWajib(respons, 'daftar mahasiswa')
  const metadata = objekWajib(data.meta, 'daftar mahasiswa.meta')

  if (!Array.isArray(data.data)) {
    throw new Error('Respons API daftar mahasiswa tidak memuat array data')
  }

  const pagination = {
    halaman: Number(metadata.current_page ?? metadata.halaman),
    perHalaman: Number(metadata.per_page ?? metadata.per_halaman),
    total: Number(metadata.total),
    totalHalaman: Number(metadata.last_page ?? metadata.total_halaman),
  }

  if (Object.values(pagination).some((nilai) => !Number.isFinite(nilai))) {
    throw new Error('Respons API daftar mahasiswa memiliki metadata pagination yang tidak valid')
  }

  return {
    data: data.data,
    pagination,
  }
}

export function normalisasiDetailMahasiswa(respons) {
  return objekWajib(objekWajib(respons, 'detail mahasiswa').data, 'detail mahasiswa.data')
}

export function normalisasiGalatApi(error) {
  const data = error?.response?.data
  const fieldErrors = {}

  if (Array.isArray(data?.galat)) {
    for (const galat of data.galat) {
      if (typeof galat.field !== 'string' || typeof galat.message !== 'string') {
        throw new Error('Format galat validasi Fiber tidak valid')
      }

      fieldErrors[galat.field] ??= []
      fieldErrors[galat.field].push(galat.message)
    }
  } else if (data?.errors && typeof data.errors === 'object') {
    for (const [field, messages] of Object.entries(data.errors)) {
      fieldErrors[field] = Array.isArray(messages) ? messages : [String(messages)]
    }
  }

  return {
    status: error?.response?.status,
    message: data?.pesan ?? data?.message ?? error?.message ?? 'Permintaan gagal',
    fieldErrors,
  }
}
