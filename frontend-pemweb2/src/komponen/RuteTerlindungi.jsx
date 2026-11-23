import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { KUNCI_TOKEN } from '../api/klien.js'

function RuteTerlindungi() {
  const lokasi = useLocation()
  const token = sessionStorage.getItem(KUNCI_TOKEN)

  if (!token) {
    return <Navigate to="/login" replace state={{ from: lokasi }} />
  }

  return <Outlet />
}

export default RuteTerlindungi
