import axios from 'axios'
import { konfigurasiApi } from './kontrak.js'

export const KUNCI_TOKEN = 'token'

export const klien = axios.create({
  baseURL: konfigurasiApi.baseURL,
  timeout: 10000,
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json',
  },
})

klien.interceptors.request.use((config) => {
  const token = sessionStorage.getItem(KUNCI_TOKEN)

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  return config
})

klien.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error?.response?.status === 401) {
      sessionStorage.removeItem(KUNCI_TOKEN)

      if (window.location.pathname !== '/login') {
        window.location.assign('/login')
      }
    }

    return Promise.reject(error)
  },
)
