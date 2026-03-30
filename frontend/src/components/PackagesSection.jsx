const formatDuration = (pkg) => {
  if (pkg?.duration_minutes) {
    if (pkg.duration_minutes % 1440 === 0) {
      const days = pkg.duration_minutes / 1440
      return `${days} day${days === 1 ? '' : 's'}`
    }
    if (pkg.duration_minutes % 60 === 0) {
      const hours = pkg.duration_minutes / 60
      return `${hours} hour${hours === 1 ? '' : 's'}`
    }
    return `${pkg.duration_minutes} minutes`
  }
  if (pkg?.validity_days) {
    return `${pkg.validity_days} day${pkg.validity_days === 1 ? '' : 's'}`
  }
  return 'Flexible'
}

export const PackagesSection = ({
  packages,
  packagesLoading,
  packagesError,
  onReload,
}) => {
  return (
    <section id="packages" className="section">
      <div className="section-title">
        <div>
          <h3>Packages</h3>
          <p>Live prices pulled from your backend.</p>
        </div>
        <button className="btn ghost" onClick={onReload}>
          Reload
        </button>
      </div>
      {packagesLoading ? (
        <div className="empty">Loading packages...</div>
      ) : packagesError ? (
        <div className="empty error">{packagesError}</div>
      ) : (
        <div className="grid">
          {packages.map((pkg) => (
            <article className="card" key={pkg.id}>
              <p className="tag">{pkg.package_type?.toUpperCase()}</p>
              <h4>{pkg.name}</h4>
              <p className="muted">{pkg.description || 'Unlimited access'}</p>
              <div className="price-row">
                <span className="price">Ksh {pkg.price}</span>
                <span className="duration">{formatDuration(pkg)}</span>
              </div>
            </article>
          ))}
        </div>
      )}
    </section>
  )
}
