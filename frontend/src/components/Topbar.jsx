export const Topbar = ({ apiBase, packagesLoading, packagesCount }) => {
  return (
    <header className="topbar">
      <div className="brand">
        <div className="brand-icon" aria-hidden="true">
          <span className="wifi-dot" />
          <span className="wifi-arc" />
          <span className="wifi-arc" />
        </div>
        <div className="brand-copy">
          <h1 className="brand-title">CEATON</h1>
          <p className="brand-subtitle">Ultra-Fast WiFi Experience</p>
          <p className="brand-help">
            <span className="support-label">Support:</span> +254 795 671 930
          </p>
        </div>
      </div>
      <div className="status-card" aria-hidden="true">
        <p className="label">API Base</p>
        <p className="value">{apiBase}</p>
        <p className="status">
          {packagesLoading ? 'Loading packages' : `${packagesCount} packages`}
        </p>
      </div>
    </header>
  )
}
