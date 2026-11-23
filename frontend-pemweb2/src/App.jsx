import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { KUNCI_TOKEN } from './api/klien.js'
import DaftarMahasiswa from './halaman/DaftarMahasiswa.jsx'
import DemoMahasiswa from './halaman/DemoMahasiswa.jsx'
import Login from './halaman/Login.jsx'
import RuteTerlindungi from './komponen/RuteTerlindungi.jsx'

function Beranda() {
  const tujuan = sessionStorage.getItem(KUNCI_TOKEN) ? '/mahasiswa' : '/login'
  return <Navigate to={tujuan} replace />
}

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/demo-mahasiswa" element={<DemoMahasiswa />} />
        <Route element={<RuteTerlindungi />}>
          <Route path="/mahasiswa" element={<DaftarMahasiswa />} />
        </Route>
        <Route path="/" element={<Beranda />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
