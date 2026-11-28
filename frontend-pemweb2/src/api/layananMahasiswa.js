import { klien } from './klien.js'
import {
  buatParameterDaftarMahasiswa,
  normalisasiDaftarMahasiswa,
  normalisasiDetailMahasiswa,
} from './kontrak.js'

export async function ambilDaftarMahasiswa({
  halaman = 1,
  perHalaman = 10,
  cari,
  angkatan,
} = {}) {
  const { data } = await klien.get('/mahasiswa', {
    params: buatParameterDaftarMahasiswa({
      halaman,
      perHalaman,
      cari,
      angkatan,
    }),
  })

  return normalisasiDaftarMahasiswa(data)
}

export async function ambilDetailMahasiswa(id) {
  const { data } = await klien.get(`/mahasiswa/${encodeURIComponent(id)}`)
  return normalisasiDetailMahasiswa(data)
}

export async function tambahMahasiswa(mahasiswa) {
  const { data } = await klien.post('/mahasiswa', mahasiswa)
  return normalisasiDetailMahasiswa(data)
}

export async function perbaruiMahasiswa(id, mahasiswa) {
  const { data } = await klien.put(`/mahasiswa/${encodeURIComponent(id)}`, mahasiswa)
  return normalisasiDetailMahasiswa(data)
}

export async function hapusMahasiswa(id) {
  const { data } = await klien.delete(`/mahasiswa/${encodeURIComponent(id)}`)
  return data
}
