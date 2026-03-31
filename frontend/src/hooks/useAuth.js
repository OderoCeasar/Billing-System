import { useEffect, useState } from 'react'

const TOKEN_KEY = 'auth_token'

export const useAuth = (apiFetch) => {
  const [token, setToken] = useState(() => localStorage.getItem(TOKEN_KEY) || '')
  const [user, setUser] = useState(null)
  const [authMessage, setAuthMessage] = useState('')

  useEffect(() => {
    if (token) {
      localStorage.setItem(TOKEN_KEY, token)
    } else {
      localStorage.removeItem(TOKEN_KEY)
    }
  }, [token])

  const handleAuthSuccess = (newToken, newUser) => {
    setToken(newToken)
    setUser(newUser || null)
    setAuthMessage('Authenticated successfully')
  }

  const login = async (payload) => {
    setAuthMessage('')
    const data = await apiFetch('/auth/login', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    handleAuthSuccess(data.token, data.user)
    return data
  }

  const register = async (payload) => {
    setAuthMessage('')
    const data = await apiFetch('/auth/register', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    handleAuthSuccess(data.token, data.user)
    return data
  }

  const quickRegister = async (phoneNumber) => {
    setAuthMessage('')
    const data = await apiFetch('/auth/quick-register', {
      method: 'POST',
      body: JSON.stringify({ phone_number: phoneNumber }),
    })
    handleAuthSuccess(data.token, data.user)
    return data
  }

  const logout = () => {
    setToken('')
    setUser(null)
    setAuthMessage('Logged out')
  }

  return {
    token,
    user,
    authMessage,
    setAuthMessage,
    login,
    register,
    quickRegister,
    logout,
  }
}
