import { useState } from 'react'
import './DaftarMahasiswa.css'

const DAFTAR_AWAL = [
  {
    id: 1,
    nim: 'H2025001',
    nama: 'Contoh Mahasiswa Satu',
    angkatan: 2025,
    ipk: 3.56,
  },
  {
    id: 2,
    nim: 'H2025002',
    nama: 'Contoh Mahasiswa Dua',
    angkatan: 2025,
    ipk: 3.72,
  },
]

const FORM_AWAL = {
  program_studi_id: '1',
  nim: 'X1234',
  nama: 'Contoh Mahasiswa Baru',
  email: 'mahasiswa.demo@example.invalid',
  angkatan: '2025',
  ipk: '3.50',
}

function validasiForm(form) {
  const galat = {}
  const programStudiID = Number(form.program_studi_id)
  const angkatan = Number(form.angkatan)
  const ipk = Number(form.ipk)

  if (!Number.isInteger(programStudiID) || programStudiID < 1) {
    galat.program_studi_id = 'ID program studi harus berupa bilangan positif'
  }
  if (!/^H\d{4,19}$/.test(form.nim.trim())) {
    galat.nim = 'Format NIM tidak valid. NIM harus diawali H dan diikuti angka.'
  }
  if (form.nama.trim().length < 3 || form.nama.trim().length > 100) {
    galat.nama = 'Nama harus terdiri dari 3 sampai 100 karakter'
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email.trim()) || form.email.length > 100) {
    galat.email = 'Format email tidak valid'
  }
  if (!Number.isInteger(angkatan) || angkatan < 2000 || angkatan > 2100) {
    galat.angkatan = 'Angkatan harus berada pada rentang 2000 sampai 2100'
  }
  if (form.ipk.trim() === '' || !Number.isFinite(ipk) || ipk < 0 || ipk > 4) {
    galat.ipk = 'IPK harus berada pada rentang 0 sampai 4'
  }

  return galat
}

function DemoMahasiswa() {
  const [mahasiswa, setMahasiswa] = useState(DAFTAR_AWAL)
  const [formMahasiswa, setFormMahasiswa] = useState(FORM_AWAL)
  const [galatField, setGalatField] = useState({})
  const [pesanTambah, setPesanTambah] = useState('')
  const [pesanValidasi, setPesanValidasi] = useState('')
  const [pesanAksi, setPesanAksi] = useState('')

  function ubahForm(field, value) {
    setFormMahasiswa((sebelumnya) => ({ ...sebelumnya, [field]: value }))
    setGalatField((sebelumnya) => ({ ...sebelumnya, [field]: undefined }))
    setPesanTambah('')
    setPesanValidasi('')
  }

  function tanganiSubmit(event) {
    event.preventDefault()
    setPesanAksi('')

    const galat = validasiForm(formMahasiswa)
    setGalatField(galat)

    if (Object.keys(galat).length > 0) {
      setPesanTambah('Data belum valid. Periksa field yang ditandai.')
      return
    }

    setPesanValidasi('Data valid. Mode demo tidak mengirim atau menyimpan data.')
  }

  function tanganiHapus(data) {
    if (!window.confirm(`Hapus data mahasiswa ${data.nama} (${data.nim})?`)) {
      return
    }

    setMahasiswa((sebelumnya) => sebelumnya.filter((item) => item.id !== data.id))
    setGalatField({})
    setPesanTambah('')
    setPesanValidasi('')
    setPesanAksi(`Data mahasiswa ${data.nama} berhasil dihapus dari daftar demo.`)
  }

  function pesanField(field) {
    return galatField[field]
  }

  return (
    <main className="mahasiswa-page">
      <header className="mahasiswa-header">
        <div>
          <p className="mahasiswa-eyebrow">DEMO FRONTEND — DATA CONTOH</p>
          <h1>Daftar Mahasiswa</h1>
          <p className="mahasiswa-description">
            Demo lokal untuk screenshot; tidak terhubung ke API atau database.
          </p>
        </div>
        <p className="mahasiswa-total">{mahasiswa.length} mahasiswa demo</p>
      </header>

      <section className="mahasiswa-panel" aria-label="Demo daftar mahasiswa">
        <details className="mahasiswa-add" open>
          <summary>Tambah mahasiswa</summary>
          <form className="mahasiswa-form" noValidate onSubmit={tanganiSubmit}>
            <div className="mahasiswa-form-grid">
              <div className="mahasiswa-form-field">
                <label htmlFor="demo-program-studi">ID Program Studi</label>
                <input
                  id="demo-program-studi"
                  type="number"
                  min="1"
                  step="1"
                  value={formMahasiswa.program_studi_id}
                  onChange={(event) => ubahForm('program_studi_id', event.target.value)}
                  aria-invalid={Boolean(pesanField('program_studi_id'))}
                  aria-describedby={pesanField('program_studi_id') ? 'demo-galat-program-studi' : undefined}
                  required
                />
                {pesanField('program_studi_id') && (
                  <p className="mahasiswa-field-error" id="demo-galat-program-studi">
                    {pesanField('program_studi_id')}
                  </p>
                )}
              </div>

              <div className="mahasiswa-form-field">
                <label htmlFor="demo-nim">NIM</label>
                <input
                  id="demo-nim"
                  type="text"
                  minLength="5"
                  maxLength="20"
                  value={formMahasiswa.nim}
                  onChange={(event) => ubahForm('nim', event.target.value)}
                  aria-invalid={Boolean(pesanField('nim'))}
                  aria-describedby={pesanField('nim') ? 'demo-galat-nim' : undefined}
                  required
                />
                {pesanField('nim') && (
                  <p className="mahasiswa-field-error" id="demo-galat-nim">
                    {pesanField('nim')}
                  </p>
                )}
              </div>

              <div className="mahasiswa-form-field">
                <label htmlFor="demo-nama">Nama</label>
                <input
                  id="demo-nama"
                  type="text"
                  minLength="3"
                  maxLength="100"
                  value={formMahasiswa.nama}
                  onChange={(event) => ubahForm('nama', event.target.value)}
                  aria-invalid={Boolean(pesanField('nama'))}
                  aria-describedby={pesanField('nama') ? 'demo-galat-nama' : undefined}
                  required
                />
                {pesanField('nama') && (
                  <p className="mahasiswa-field-error" id="demo-galat-nama">
                    {pesanField('nama')}
                  </p>
                )}
              </div>

              <div className="mahasiswa-form-field">
                <label htmlFor="demo-email">Email</label>
                <input
                  id="demo-email"
                  type="email"
                  maxLength="100"
                  value={formMahasiswa.email}
                  onChange={(event) => ubahForm('email', event.target.value)}
                  aria-invalid={Boolean(pesanField('email'))}
                  aria-describedby={pesanField('email') ? 'demo-galat-email' : undefined}
                  required
                />
                {pesanField('email') && (
                  <p className="mahasiswa-field-error" id="demo-galat-email">
                    {pesanField('email')}
                  </p>
                )}
              </div>

              <div className="mahasiswa-form-field">
                <label htmlFor="demo-angkatan">Angkatan</label>
                <input
                  id="demo-angkatan"
                  type="number"
                  min="2000"
                  max="2100"
                  step="1"
                  value={formMahasiswa.angkatan}
                  onChange={(event) => ubahForm('angkatan', event.target.value)}
                  aria-invalid={Boolean(pesanField('angkatan'))}
                  aria-describedby={pesanField('angkatan') ? 'demo-galat-angkatan' : undefined}
                  required
                />
                {pesanField('angkatan') && (
                  <p className="mahasiswa-field-error" id="demo-galat-angkatan">
                    {pesanField('angkatan')}
                  </p>
                )}
              </div>

              <div className="mahasiswa-form-field">
                <label htmlFor="demo-ipk">IPK</label>
                <input
                  id="demo-ipk"
                  type="number"
                  min="0"
                  max="4"
                  step="0.01"
                  value={formMahasiswa.ipk}
                  onChange={(event) => ubahForm('ipk', event.target.value)}
                  aria-invalid={Boolean(pesanField('ipk'))}
                  aria-describedby={pesanField('ipk') ? 'demo-galat-ipk' : undefined}
                  required
                />
                {pesanField('ipk') && (
                  <p className="mahasiswa-field-error" id="demo-galat-ipk">
                    {pesanField('ipk')}
                  </p>
                )}
              </div>
            </div>

            <p className="mahasiswa-form-hint">
              Data pada halaman ini hanya contoh dan tidak dikirim ke backend.
            </p>
            {pesanTambah && (
              <p className="mahasiswa-form-error" role="alert">
                {pesanTambah}
              </p>
            )}
            {pesanValidasi && (
              <p className="mahasiswa-success" role="status">
                {pesanValidasi}
              </p>
            )}
            <button className="mahasiswa-submit" type="submit">
              Simpan mahasiswa
            </button>
          </form>
        </details>

        {pesanAksi && (
          <p className="mahasiswa-success" role="status">
            {pesanAksi}
          </p>
        )}

        {mahasiswa.length === 0 ? (
          <p className="mahasiswa-empty">Tidak ada data mahasiswa demo.</p>
        ) : (
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
                    <td>{data.ipk.toFixed(2)}</td>
                    <td>
                      <button
                        className="mahasiswa-delete"
                        type="button"
                        onClick={() => tanganiHapus(data)}
                      >
                        Hapus
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </main>
  )
}

export default DemoMahasiswa
