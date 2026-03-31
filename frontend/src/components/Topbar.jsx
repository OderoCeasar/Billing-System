export const Topbar = ({ apiBase, packagesLoading, packagesCount }) => {
  return (
    <header className="topbar">
      <div className="brand">
        <span className="brand-dot" />
        <div>
          <p className="eyebrow">WiFi Billing System</p>
          <h1>Smart access for every device.</h1>
        </div>
      </div>
      <div className="status-card">
        <p className="label">API Base</p>
        <p className="value">{apiBase}</p>
        <p className="status">
          {packagesLoading ? 'Loading packages' : `${packagesCount} packages`}
        </p>
      </div>
    </header>
  )
}
