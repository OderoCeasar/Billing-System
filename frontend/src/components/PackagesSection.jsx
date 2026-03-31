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

const getSpeedLabel = (pkg) => {
  const up = Number(pkg?.speed_limit_up || 0)
  const down = Number(pkg?.speed_limit_down || 0)
  const peak = Math.max(up, down)
  if (peak > 0) {
    return `Up to ${peak} Mbps`
  }
  return 'High-speed access'
}

export const PackagesSection = ({
  packages,
  packagesLoading,
  packagesError,
  onSelectPackage,
}) => {
  return (
    <section id="packages" className="plans-grid">
      {packagesLoading ? (
        <div className="empty">Loading packages...</div>
      ) : packagesError ? (
        <div className="empty error">{packagesError}</div>
      ) : (
        <div className="plans">
          {packages.map((pkg, index) => {
            const isPopular = index === 1
            return (
              <article
                className={`plan-card${isPopular ? ' popular' : ''}`}
                key={pkg.id}
              >
                {isPopular ? <span className="plan-badge">POPULAR</span> : null}
                <div className="plan-header">
                  <h4>{pkg.name}</h4>
                  <span className="plan-price">
                    <span className="plan-currency">Ksh</span>
                    <span className="plan-amount">{pkg.price}</span>
                  </span>
                </div>
                <p className="plan-description">
                  {pkg.description || 'Unlimited access'}
                </p>
                <ul className="plan-features">
                  <li>
                    <span className="feature-icon" aria-hidden="true" />
                    {formatDuration(pkg)} access
                  </li>
                  <li>
                    <span className="feature-icon speed" aria-hidden="true" />
                    {getSpeedLabel(pkg)}
                  </li>
                </ul>
                <button
                  className={`btn ${isPopular ? 'primary' : 'ghost'} plan-action`}
                  type="button"
                  onClick={() => onSelectPackage?.(pkg)}
                >
                  Buy
                </button>
              </article>
            )
          })}
        </div>
      )}
    </section>
  )
}
