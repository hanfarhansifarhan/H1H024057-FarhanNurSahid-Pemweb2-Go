import { useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { normalisasiGalatApi } from '../api/kontrak.js'
import { login } from '../api/layananAuth.js'
import './Login.css'

function pesanField(galatField, field) {
  return galatField[field]?.[0]
}

function Login() {
  const lokasi = useLocation()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [kataSandi, setKataSandi] = useState('')
  const [sedangMemuat, setSedangMemuat] = useState(false)
  const [pesanGalat, setPesanGalat] = useState(
    lokasi.state?.pesanGalatLogout ?? '',
  )
  const [galatField, setGalatField] = useState({})

  async function tanganiSubmit(event) {
    event.preventDefault()
    setSedangMemuat(true)
    setPesanGalat('')
    setGalatField({})

    try {
      await login(email, kataSandi)
      navigate('/mahasiswa', { replace: true })
    } catch (error) {
      const galat = normalisasiGalatApi(error)
      setPesanGalat(galat.message)
      setGalatField(galat.fieldErrors)
    } finally {
      setSedangMemuat(false)
    }
  }

  const galatEmail = pesanField(galatField, 'email')
  const galatKataSandi =
    pesanField(galatField, 'password') ?? pesanField(galatField, 'kata_sandi')

  return (
    <main className="login-page">
      <section className="login-card" aria-labelledby="login-title">
        <p className="login-eyebrow">PEMROGRAMAN WEB II</p>
        <h1 id="login-title">Masuk</h1>
        <p className="login-description">
          Masuk untuk mengelola data mahasiswa.
        </p>

        <form className="login-form" onSubmit={tanganiSubmit}>
          <label htmlFor="email">Email</label>
          <input
            id="email"
            name="email"
            type="email"
            autoComplete="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            aria-invalid={Boolean(galatEmail)}
            aria-describedby={galatEmail ? 'email-error' : undefined}
            required
            disabled={sedangMemuat}
          />
          {galatEmail && (
            <p className="login-field-error" id="email-error">
              {galatEmail}
            </p>
          )}

          <label htmlFor="kata-sandi">Kata sandi</label>
          <input
            id="kata-sandi"
            name="kata_sandi"
            type="password"
            autoComplete="current-password"
            value={kataSandi}
            onChange={(event) => setKataSandi(event.target.value)}
            aria-invalid={Boolean(galatKataSandi)}
            aria-describedby={galatKataSandi ? 'kata-sandi-error' : undefined}
            required
            disabled={sedangMemuat}
          />
          {galatKataSandi && (
            <p className="login-field-error" id="kata-sandi-error">
              {galatKataSandi}
            </p>
          )}

          {pesanGalat && (
            <p className="login-error" role="alert">
              {pesanGalat}
            </p>
          )}

          <button type="submit" disabled={sedangMemuat}>
            {sedangMemuat ? 'Memproses...' : 'Masuk'}
          </button>
          {sedangMemuat && (
            <span className="login-loading" role="status">
              Sedang memverifikasi akun
            </span>
          )}
        </form>
      </section>
    </main>
  )
}

export default Login
