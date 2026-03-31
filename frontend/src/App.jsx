import { useMemo } from 'react'
import './App.css'
import { Topbar } from './components/Topbar'
import { PackagesSection } from './components/PackagesSection'
import { AuthSection } from './components/AuthSection'
import { VoucherSection } from './components/VoucherSection'
import { useAuth } from './hooks/useAuth'
import { usePackages } from './hooks/usePackages'
import { createApiClient } from './lib/api'

function App() {
  const apiBase =
    useMemo(() => import.meta.env.VITE_API_BASE, []) ||
    'http://localhost:8080/api/v1'
  const apiClient = useMemo(
    () => createApiClient(apiBase, () => localStorage.getItem('auth_token') || ''),
    [apiBase]
  )
  const { apiFetch } = apiClient

  const {
    token: authToken,
    authMessage,
    setAuthMessage,
    login,
    logout,
  } = useAuth(apiFetch)

  const {
    packages,
    packagesLoading,
    packagesError,
    loadPackages,
  } = usePackages(apiFetch)

  return (
    <div className="app">
      <Topbar
        apiBase={apiBase}
        packagesLoading={packagesLoading}
        packagesCount={packages.length}
      />
      <div className="cta-row">
        <AuthSection
          token={authToken}
          authMessage={authMessage}
          setAuthMessage={setAuthMessage}
          onLogin={login}
          onLogout={logout}
        />
        <VoucherSection />
      </div>
      <div className="plans-header">
        <h3>Choose Your Plan</h3>
        <p>Select the perfect package for your needs</p>
      </div>
      <PackagesSection
        packages={packages}
        packagesLoading={packagesLoading}
        packagesError={packagesError}
      />
    </div>
  )
}

export default App
