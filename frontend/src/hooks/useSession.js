import { useCallback, useEffect, useState } from 'react'

export const useSession = (apiFetch, token) => {
  const [sessionStats, setSessionStats] = useState(null)
  const [sessionMessage, setSessionMessage] = useState('')

  const loadActiveSession = useCallback(async () => {
    if (!token) {
      setSessionStats(null)
      return
    }
    try {
      const data = await apiFetch('/sessions/active', { method: 'GET' })
      setSessionStats(data)
      setSessionMessage('')
    } catch (error) {
      setSessionStats(null)
      setSessionMessage(error.message)
    }
  }, [apiFetch, token])

  useEffect(() => {
    loadActiveSession()
  }, [loadActiveSession])

  const disconnect = async () => {
    if (!sessionStats?.session?.id) {
      return
    }
    setSessionMessage('')
    await apiFetch(`/sessions/${sessionStats.session.id}/disconnect`, {
      method: 'POST',
    })
    setSessionMessage('Session disconnected')
    await loadActiveSession()
  }

  return {
    sessionStats,
    sessionMessage,
    setSessionMessage,
    loadActiveSession,
    disconnect,
  }
}
