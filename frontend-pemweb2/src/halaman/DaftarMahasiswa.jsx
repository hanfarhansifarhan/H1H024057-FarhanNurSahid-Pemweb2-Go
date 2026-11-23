import { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { normalisasiGalatApi } from '../api/kontrak.js'
import { logout } from '../api/layananAuth.js'
import {
  ambilDaftarMahasiswa,
  hapusMahasiswa,
  tambahMahasiswa,
} from '../api/layananMahasiswa.js'
import './DaftarMahasiswa.css'

const PER_HALAMAN = 10
const FORM_AWAL = {
  program_studi_id: '',
  nim: '',
  nama: '',
  email: '',
  angkatan: '',
  ipk: '',
}

function DaftarMahasiswa() {
  const navigate = useNavigate()
  const [kataCari, setKataCari] = useState('')
  const [cariAktif, setCariAktif] = useState('')
  const cariAktifRef = useRef('')
  const [halaman, setHalaman] = useState(1)
  const [mahasiswa, setMahasiswa] = useState([])
  const [pagination, setPagination] = useState({
    halaman: 1,
    perHalaman: PER_HALAMAN,
    total: 0,
    totalHalaman: 1,
  })
  const [sedangMemuat, setSedangMemuat] = useState(true)
  const [pesanGalat, setPesanGalat] = useState('')
  const [pembaruanDaftar, setPembaruanDaftar] = useState(0)
  const [sedangKeluar, setSedangKeluar] = useState(false)
  const [formMahasiswa, setFormMahasiswa] = useState(FORM_AWAL)
  const [sedangMenambah, setSedangMenambah] = useState(false)
  const [pesanTambah, setPesanTambah] = useState('')
  const [galatTambah, setGalatTambah] = useState({})
  const [pesanAksi, setPesanAksi] = useState('')
  const [idMenghapus, setIdMenghapus] = useState(null)

  useEffect(() => {
    const timer = setTimeout(() => {
      const cariBaru = kataCari.trim()
      if (cariBaru !== cariAktifRef.current) {
        cariAktifRef.current = cariBaru
        setCariAktif(cariBaru)
        setHalaman(1)
        setSedangMemuat(true)
        setPesanGalat('')
      }
    }, 400)

    return () => clearTimeout(timer)
  }, [kataCari])

  useEffect(() => {
    let masihAktif = true

    ambilDaftarMahasiswa({
      halaman,
      perHalaman: PER_HALAMAN,
      cari: cariAktif || undefined,
    })
      .then((hasil) => {
        if (masihAktif) {
          setMahasiswa(hasil.data)
          setPagination(hasil.pagination)
        }
      })
      .catch((error) => {
        if (masihAktif) {
          setPesanGalat(normalisasiGalatApi(error).message)
        }
      })
      .finally(() => {
        if (masihAktif) {
          setSedangMemuat(false)
        }
      })

    return () => {
      masihAktif = false
    }
  }, [halaman, cariAktif, pembaruanDaftar])

  const totalHalaman = Math.max(1, pagination.totalHalaman)

  function ubahHalaman(halamanBaru) {
    setPesanGalat('')
    setSedangMemuat(true)
    setHalaman(halamanBaru)
  }

  function cobaMuatUlang() {
    setPesanGalat('')
    setSedangMemuat(true)
    setPembaruanDaftar((nilai) => nilai + 1)
  }

  function ubahFormMahasiswa(field, value) {
    setFormMahasiswa((sebelumnya) => ({ ...sebelumnya, [field]: value }))
    setGalatTambah((sebelumnya) => ({ ...sebelumnya, [field]: undefined }))
    setPesanTambah('')
  }

  function galatField(field) {
    return galatTambah[field]?.[0]
  }

  async function tanganiTambah(event) {
    event.preventDefault()
    setSedangMenambah(true)
    setPesanTambah('')
    setGalatTambah({})
    setPesanAksi('')

    try {
      await tambahMahasiswa({
        program_studi_id: Number(formMahasiswa.program_studi_id),
        nim: formMahasiswa.nim,
        nama: formMahasiswa.nama,
        email: formMahasiswa.email,
        angkatan: Number(formMahasiswa.angkatan),
        ipk: Number(formMahasiswa.ipk),
      })
      setFormMahasiswa(FORM_AWAL)
      setPesanAksi('Data mahasiswa berhasil ditambahkan.')
      setPesanGalat('')
      setHalaman(1)
      setPembaruanDaftar((nilai) => nilai + 1)
    } catch (error) {
      const galat = normalisasiGalatApi(error)
      setPesanTambah(galat.message)
      setGalatTambah(galat.fieldErrors)
    } finally {
      setSedangMenambah(false)
    }
  }

  async function tanganiHapus(data) {
    if (!window.confirm(`Hapus data mahasiswa ${data.nama} (${data.nim})?`)) {
      return
    }

    setIdMenghapus(data.id)
    setPesanGalat('')
    setPesanAksi('')

    try {
      await hapusMahasiswa(data.id)
      setPesanAksi(`Data mahasiswa ${data.nama} berhasil dihapus.`)
      setHalaman(1)
      setPembaruanDaftar((nilai) => nilai + 1)
    } catch (error) {
      setPesanGalat(normalisasiGalatApi(error).message)
    } finally {
      setIdMenghapus(null)
    }
  }

  async function tanganiLogout() {
    setSedangKeluar(true)

    try {
      await logout()
      navigate('/login', { replace: true })
    } catch (error) {
      const galat = normalisasiGalatApi(error)
      navigate('/login', {
        replace: true,
        state: {
          pesanGalatLogout: `Sesi lokal sudah dihapus, tetapi server tidak mengonfirmasi logout: ${galat.message}`,
        },
      })
    }
  }

  return (
    <main className="mahasiswa-page">
      <header className="mahasiswa-header">
        <div>
          <p className="mahasiswa-eyebrow">DATA AKADEMIK</p>
          <h1>Daftar Mahasiswa</h1>
          <p className="mahasiswa-description">
            Cari dan lihat data mahasiswa.
          </p>
        </div>
        <div className="mahasiswa-header-actions">
          <p className="mahasiswa-total">{pagination.total} mahasiswa</p>
          <button
            className="mahasiswa-logout"
            type="button"
            onClick={tanganiLogout}
            disabled={sedangKeluar}
          >
            {sedangKeluar ? 'Keluar...' : 'Keluar'}
          </button>
        </div>
      </header>

      <section className="mahasiswa-panel" aria-label="Daftar mahasiswa">
        <details className="mahasiswa-add" open>
          <summary>Tambah mahasiswa</summary>
          <form className="mahasiswa-form" onSubmit={tanganiTambah}>
            <div className="mahasiswa-form-grid">
              <div className="mahasiswa-form-field">
                <label htmlFor="form-program-studi">ID Program Studi</label>
                <input
                  id="form-program-studi"
                  type="number"
                  min="1"
                  step="1"
                  value={formMahasiswa.program_studi_id}
                  onChange={(event) => ubahFormMahasiswa('program_studi_id', event.target.value)}
                  aria-invalid={Boolean(galatField('program_studi_id'))}
                  aria-describedby={galatField('program_studi_id') ? 'galat-program-studi' : undefined}
                  required
                  disabled={sedangMenambah}
                />
                {galatField('program_studi_id') && (
                  <p className="mahasiswa-field-error" id="galat-program-studi">
                    {galatField('program_studi_id')}
                  </p>
                )}
              </div>
              <div className="mahasiswa-form-field">
                <label htmlFor="form-nim">NIM</label>
                <input
                  id="form-nim"
                  type="text"
                  minLength="5"
                  maxLength="20"
                  value={formMahasiswa.nim}
                  onChange={(event) => ubahFormMahasiswa('nim', event.target.value)}
                  aria-invalid={Boolean(galatField('nim'))}
                  aria-describedby={galatField('nim') ? 'galat-nim' : undefined}
                  required
                  disabled={sedangMenambah}
                />
                {galatField('nim') && (
                  <p className="mahasiswa-field-error" id="galat-nim">
                    {galatField('nim')}
                  </p>
                )}
              </div>
              <div className="mahasiswa-form-field">
                <label htmlFor="form-nama">Nama</label>
                <input
                  id="form-nama"
                  type="text"
                  minLength="3"
                  maxLength="100"
                  value={formMahasiswa.nama}
                  onChange={(event) => ubahFormMahasiswa('nama', event.target.value)}
                  aria-invalid={Boolean(galatField('nama'))}
                  aria-describedby={galatField('nama') ? 'galat-nama' : undefined}
                  required
                  disabled={sedangMenambah}
                />
                {galatField('nama') && (
                  <p className="mahasiswa-field-error" id="galat-nama">
                    {galatField('nama')}
                  </p>
                )}
              </div>
              <div className="mahasiswa-form-field">
                <label htmlFor="form-email">Email</label>
                <input
                  id="form-email"
                  type="email"
                  maxLength="100"
                  value={formMahasiswa.email}
                  onChange={(event) => ubahFormMahasiswa('email', event.target.value)}
                  aria-invalid={Boolean(galatField('email'))}
                  aria-describedby={galatField('email') ? 'galat-email' : undefined}
                  required
                  disabled={sedangMenambah}
                />
                {galatField('email') && (
                  <p className="mahasiswa-field-error" id="galat-email">
                    {galatField('email')}
                  </p>
                )}
              </div>
              <div className="mahasiswa-form-field">
                <label htmlFor="form-angkatan">Angkatan</label>
                <input
                  id="form-angkatan"
                  type="number"
                  min="2000"
                  max="2100"
                  step="1"
                  value={formMahasiswa.angkatan}
                  onChange={(event) => ubahFormMahasiswa('angkatan', event.target.value)}
                  aria-invalid={Boolean(galatField('angkatan'))}
                  aria-describedby={galatField('angkatan') ? 'galat-angkatan' : undefined}
                  required
                  disabled={sedangMenambah}
                />
                {galatField('angkatan') && (
                  <p className="mahasiswa-field-error" id="galat-angkatan">
                    {galatField('angkatan')}
                  </p>
                )}
              </div>
              <div className="mahasiswa-form-field">
                <label htmlFor="form-ipk">IPK</label>
                <input
                  id="form-ipk"
                  type="number"
                  min="0"
                  max="4"
                  step="0.01"
                  value={formMahasiswa.ipk}
                  onChange={(event) => ubahFormMahasiswa('ipk', event.target.value)}
                  aria-invalid={Boolean(galatField('ipk'))}
                  aria-describedby={galatField('ipk') ? 'galat-ipk' : undefined}
                  required
                  disabled={sedangMenambah}
                />
                {galatField('ipk') && (
                  <p className="mahasiswa-field-error" id="galat-ipk">
                    {galatField('ipk')}
                  </p>
                )}
              </div>
            </div>

            <p className="mahasiswa-form-hint">
              Gunakan ID program studi yang sudah terdaftar.
            </p>
            {pesanTambah && (
              <p className="mahasiswa-form-error" role="alert">
                {pesanTambah}
              </p>
            )}
            <button
              className="mahasiswa-submit"
              type="submit"
              disabled={sedangMenambah}
            >
              {sedangMenambah ? 'Menyimpan...' : 'Simpan mahasiswa'}
            </button>
          </form>
        </details>

        <label className="mahasiswa-search-label" htmlFor="cari-mahasiswa">
          Pencarian
        </label>
        <input
          id="cari-mahasiswa"
          className="mahasiswa-search"
          type="search"
          value={kataCari}
          onChange={(event) => setKataCari(event.target.value)}
          placeholder="Cari nama atau NIM..."
          aria-label="Cari berdasarkan nama atau NIM"
        />

        {pesanGalat && (
          <div className="mahasiswa-error" role="alert">
            <p>{pesanGalat}</p>
            <button type="button" onClick={cobaMuatUlang}>
              Coba lagi
            </button>
          </div>
        )}
        {pesanAksi && (
          <p className="mahasiswa-success" role="status">
            {pesanAksi}
          </p>
        )}

        {sedangMemuat ? (
          <p className="mahasiswa-status" role="status">
            Memuat data mahasiswa...
          </p>
        ) : pesanGalat ? null : mahasiswa.length === 0 ? (
          <p className="mahasiswa-empty">Data mahasiswa tidak ditemukan.</p>
        ) : (
          <>
            <div className="mahasiswa-table-wrap">
              <table className="mahasiswa-table">
                <thead>
                  <tr>
                    <th scope="col">NIM</th>
                    <th scope="col">Nama</th>
                    <th scope="col">Angkatan</th>
                    <th scope="col">IPK</th>
                    <th scope="col">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {mahasiswa.map((data) => (
                    <tr key={data.id}>
                      <td>{data.nim}</td>
                      <td>{data.nama}</td>
                      <td>{data.angkatan}</td>
                      <td>{Number(data.ipk).toFixed(2)}</td>
                      <td>
                        <button
                          className="mahasiswa-delete"
                          type="button"
                          onClick={() => tanganiHapus(data)}
                          disabled={idMenghapus !== null}
                        >
                          {idMenghapus === data.id ? 'Menghapus...' : 'Hapus'}
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <nav className="mahasiswa-pagination" aria-label="Pagination">
              <button
                type="button"
                onClick={() => ubahHalaman(Math.max(1, halaman - 1))}
                disabled={sedangMemuat || halaman <= 1}
              >
                Sebelumnya
              </button>
              <span>
                Halaman {halaman} dari {totalHalaman}
              </span>
              <button
                type="button"
                onClick={() => ubahHalaman(Math.min(totalHalaman, halaman + 1))}
                disabled={sedangMemuat || halaman >= totalHalaman}
              >
                Berikutnya
              </button>
            </nav>
          </>
        )}
      </section>
    </main>
  )
}

export default DaftarMahasiswa
