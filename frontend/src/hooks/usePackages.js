import { useCallback, useEffect, useState } from 'react'

export const usePackages = (apiFetch) => {
  const [packages, setPackages] = useState([])
  const [packagesLoading, setPackagesLoading] = useState(true)
  const [packagesError, setPackagesError] = useState('')

  const loadPackages = useCallback(async () => {
    setPackagesLoading(true)
    setPackagesError('')
    try {
      const data = await apiFetch('/packages', { method: 'GET' })
      setPackages(data.packages || [])
    } catch (error) {
      setPackagesError(error.message)
    } finally {
      setPackagesLoading(false)
    }
  }, [apiFetch])

  useEffect(() => {
    loadPackages()
  }, [loadPackages])

  return {
    packages,
    packagesLoading,
    packagesError,
    loadPackages,
  }
}
