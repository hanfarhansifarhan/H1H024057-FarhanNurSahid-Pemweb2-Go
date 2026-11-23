import { KUNCI_TOKEN, klien } from './klien.js'
import {
  buatPayloadLogin,
  normalisasiLogin,
  normalisasiProfil,
} from './kontrak.js'

export async function login(email, kataSandi) {
  const { data } = await klien.post(
    '/auth/login',
    buatPayloadLogin(email, kataSandi),
  )
  const hasil = normalisasiLogin(data)

  sessionStorage.setItem(KUNCI_TOKEN, hasil.token)
  return hasil
}

export async function ambilProfil() {
  const { data } = await klien.get('/auth/profil')
  return normalisasiProfil(data)
}

export async function logout() {
  try {
    const { data } = await klien.post('/auth/logout')
    return data
  } finally {
    sessionStorage.removeItem(KUNCI_TOKEN)
  }
}
