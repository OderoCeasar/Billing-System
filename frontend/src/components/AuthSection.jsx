import { useState } from 'react'

export const AuthSection = ({
  token,
  authMessage,
  setAuthMessage,
  onLogin,
  onRegister,
  onQuickRegister,
  onLogout,
}) => {
  const [loginForm, setLoginForm] = useState({
    phone_number: '',
    password: '',
  })
  const [registerForm, setRegisterForm] = useState({
    phone_number: '',
    password: '',
  })
  const [quickPhone, setQuickPhone] = useState('')

  const handleLogin = async (event) => {
    event.preventDefault()
    try {
      await onLogin(loginForm)
    } catch (error) {
      setAuthMessage(error.message)
    }
  }

  const handleRegister = async (event) => {
    event.preventDefault()
    try {
      await onRegister(registerForm)
    } catch (error) {
      setAuthMessage(error.message)
    }
  }

  const handleQuick = async (event) => {
    event.preventDefault()
    try {
      await onQuickRegister(quickPhone)
    } catch (error) {
      setAuthMessage(error.message)
    }
  }

  return (
    <section className="section" id="auth">
      <div className="section-title">
        <div>
          <h3>Authentication</h3>
          <p>Register, login, or quick-register via phone number.</p>
        </div>
        {token ? (
          <button className="btn ghost" onClick={onLogout}>
            Logout
          </button>
        ) : null}
      </div>
      <div className="columns">
        <form className="panel" onSubmit={handleLogin}>
          <h4>Login</h4>
          <label>
            Phone Number
            <input
              type="tel"
              value={loginForm.phone_number}
              onChange={(e) =>
                setLoginForm((prev) => ({
                  ...prev,
                  phone_number: e.target.value,
                }))
              }
              placeholder="07XXXXXXXX"
            />
          </label>
          <label>
            Password
            <input
              type="password"
              value={loginForm.password}
              onChange={(e) =>
                setLoginForm((prev) => ({
                  ...prev,
                  password: e.target.value,
                }))
              }
              placeholder="********"
            />
          </label>
          <button className="btn primary" type="submit">
            Login
          </button>
        </form>

        <form className="panel" onSubmit={handleRegister}>
          <h4>Register</h4>
          <label>
            Phone Number
            <input
              type="tel"
              value={registerForm.phone_number}
              onChange={(e) =>
                setRegisterForm((prev) => ({
                  ...prev,
                  phone_number: e.target.value,
                }))
              }
              placeholder="07XXXXXXXX"
            />
          </label>
          <label>
            Password
            <input
              type="password"
              value={registerForm.password}
              onChange={(e) =>
                setRegisterForm((prev) => ({
                  ...prev,
                  password: e.target.value,
                }))
              }
              placeholder="Create password"
            />
          </label>
          <button className="btn secondary" type="submit">
            Create Account
          </button>
        </form>

        <form className="panel" onSubmit={handleQuick}>
          <h4>Quick Register</h4>
          <label>
            Phone Number
            <input
              type="tel"
              value={quickPhone}
              onChange={(e) => setQuickPhone(e.target.value)}
              placeholder="07XXXXXXXX"
            />
          </label>
          <button className="btn ghost" type="submit">
            Generate Account
          </button>
          <p className="muted small">
            Creates a user with a random password and returns a token.
          </p>
        </form>
      </div>
      {authMessage ? <div className="notice">{authMessage}</div> : null}
    </section>
  )
}
