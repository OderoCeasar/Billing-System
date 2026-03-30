export const createApiClient = (baseUrl, getToken) => {
  const apiFetch = async (path, options = {}) => {
    const headers = {
      'Content-Type': 'application/json',
      ...(options.headers || {}),
    }

    const token = getToken?.()
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }

    const response = await fetch(`${baseUrl}${path}`, {
      ...options,
      headers,
    })

    const data = await response.json().catch(() => ({}))
    if (!response.ok) {
      const message = data?.error || 'Request failed'
      throw new Error(message)
    }

    return data
  }

  return { apiFetch }
}
