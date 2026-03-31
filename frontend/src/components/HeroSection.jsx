export const HeroSection = ({ token, sessionStats, onRefresh }) => {
  const statusLabel = sessionStats?.session
    ? sessionStats.session.status
    : token
    ? 'No active session'
    : 'Login to check'

  const remainingLabel = sessionStats?.stats?.remaining_minutes
    ? `${sessionStats.stats.remaining_minutes} mins left`
    : sessionStats?.stats?.remaining_data_mb
    ? `${sessionStats.stats.remaining_data_mb.toFixed(1)} MB left`
    : 'Session info will appear here'

  return (
    <section className="hero">
      <div>
        <h2>Sell time, not headaches.</h2>
        <p className="lead">
          Automate WiFi access, payments, and sessions in one place. The
          frontend is already wired to your API endpoints so you can test
          end-to-end flows immediately.
        </p>
        <div className="actions">
          <a className="btn primary" href="#packages">
            View Packages
          </a>
          <a className="btn ghost" href="#payments">
            Start Payment
          </a>
        </div>
        <div className="meta">
          <div>
            <span className="meta-label">Payments</span>
            <span className="meta-value">M-Pesa STK push</span>
          </div>
          <div>
            <span className="meta-label">Sessions</span>
            <span className="meta-value">Auto expiry + RADIUS</span>
          </div>
        </div>
      </div>
      <div className="hero-card">
        <p className="label">Active Session</p>
        <p className="value">{statusLabel}</p>
        <p className="sub">{remainingLabel}</p>
        <button className="btn secondary" onClick={onRefresh}>
          Refresh Session
        </button>
      </div>
    </section>
  )
}
