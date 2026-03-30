import { useEffect, useMemo, useRef } from 'react'
import './App.css'
import { Topbar } from './components/Topbar'
import { HeroSection } from './components/HeroSection'
import { PackagesSection } from './components/PackagesSection'
import { AuthSection } from './components/AuthSection'
import { PaymentsSection } from './components/PaymentsSection'
import { SessionsSection } from './components/SessionsSection'
import { useAuth } from './hooks/useAuth'
import { usePackages } from './hooks/usePackages'
import { usePayments } from './hooks/usePayments'
import { useSession } from './hooks/useSession'
import { createApiClient } from './lib/api'

function App() {
  const apiBase =
    useMemo(() => import.meta.env.VITE_API_BASE, []) ||
    'http://localhost:8080/api/v1'
  const tokenRef = useRef(localStorage.getItem('auth_token') || '')
  const { apiFetch } = createApiClient(apiBase, () => tokenRef.current)

  const {
    token: authToken,
    authMessage,
    setAuthMessage,
    login,
    register,
    quickRegister,
    logout,
  } = useAuth(apiFetch)

  const {
    packages,
    packagesLoading,
    packagesError,
    loadPackages,
  } = usePackages(apiFetch)

  const {
    paymentForm,
    setPaymentForm,
    paymentMessage,
    setPaymentMessage,
    initiatePayment,
  } = usePayments(apiFetch)

  const {
    sessionStats,
    sessionMessage,
    setSessionMessage,
    loadActiveSession,
    disconnect,
  } = useSession(apiFetch, authToken)

  useEffect(() => {
    tokenRef.current = authToken
  }, [authToken])

  return (
    <div className="app">
      <Topbar
        apiBase={apiBase}
        packagesLoading={packagesLoading}
        packagesCount={packages.length}
      />
      <HeroSection
        token={authToken}
        sessionStats={sessionStats}
        onRefresh={loadActiveSession}
      />
      <PackagesSection
        packages={packages}
        packagesLoading={packagesLoading}
        packagesError={packagesError}
        onReload={loadPackages}
      />
      <AuthSection
        token={authToken}
        authMessage={authMessage}
        setAuthMessage={setAuthMessage}
        onLogin={login}
        onRegister={register}
        onQuickRegister={quickRegister}
        onLogout={logout}
      />
      <PaymentsSection
        packages={packages}
        token={authToken}
        paymentForm={paymentForm}
        setPaymentForm={setPaymentForm}
        paymentMessage={paymentMessage}
        setPaymentMessage={setPaymentMessage}
        onSubmit={async () => {
          await initiatePayment()
          await loadActiveSession()
        }}
      />
      <SessionsSection
        sessionStats={sessionStats}
        sessionMessage={sessionMessage}
        setSessionMessage={setSessionMessage}
        onRefresh={loadActiveSession}
        onDisconnect={disconnect}
      />
    </div>
  )
}

export default App
