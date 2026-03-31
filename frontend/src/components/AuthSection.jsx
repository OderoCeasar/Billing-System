import { useState } from 'react'

export const AuthSection = ({
  token,
  authMessage,
  setAuthMessage,
  onLogin,
  onLogout,
}) => {
  const [isOpen, setIsOpen] = useState(false)
  const [loginForm, setLoginForm] = useState({ phone_number: '', password: '' })

  const handleLogin = async (event) => {
    event.preventDefault()
    try {
      await onLogin(loginForm)
      setIsOpen(false)
    } catch (error) {
      setAuthMessage(error.message)
    }
  }

  return (
    <section
      className={`cta-card login-card${token ? ' is-authenticated' : ''}`}
      id="auth"
      role="button"
      tabIndex={0}
      onClick={() => setIsOpen(true)}
      onKeyDown={(event) => {
        if (event.key === 'Enter' || event.key === ' ') {
          event.preventDefault()
          setIsOpen(true)
        }
      }}
    >
      <div className="cta-icon shield" aria-hidden="true" />
      <h3>Login</h3>
      <p>Access your account</p>
      {token ? (
        <button
          className="btn ghost logout-btn"
          onClick={(event) => {
            event.stopPropagation()
            onLogout()
          }}
          type="button"
        >
          Logout
        </button>
      ) : null}
      {authMessage ? <div className="notice">{authMessage}</div> : null}

      {isOpen ? (
        <div className="modal-backdrop" onClick={() => setIsOpen(false)}>
          <div className="modal" onClick={(event) => event.stopPropagation()}>
            <div className="modal-header">
              <div>
                <p className="modal-title">Access Your Account</p>
              </div>
              <button className="modal-close" type="button" onClick={() => setIsOpen(false)}>
                ×
              </button>
            </div>
            <div className="modal-tabs" role="tablist" aria-label="Login methods">
              <button className="tab active" type="button">
                Username
              </button>
              <button className="tab" type="button">
                Voucher
              </button>
              <button className="tab" type="button">
                Receipt
              </button>
            </div>
            <form className="modal-body" onSubmit={handleLogin}>
              <label>
                Username
                <input
                  type="tel"
                  autoComplete="username"
                  value={loginForm.phone_number}
                  onChange={(e) =>
                    setLoginForm((prev) => ({
                      ...prev,
                      phone_number: e.target.value,
                    }))
                  }
                  placeholder="Enter your username"
                />
              </label>
              <label>
                Password
                <input
                  type="password"
                  autoComplete="current-password"
                  value={loginForm.password}
                  onChange={(e) =>
                    setLoginForm((prev) => ({
                      ...prev,
                      password: e.target.value,
                    }))
                  }
                  placeholder="Enter your password"
                />
              </label>
              <button className="btn primary modal-submit" type="submit">
                Login Now
              </button>
            </form>
          </div>
        </div>
      ) : null}
    </section>
  )
}
