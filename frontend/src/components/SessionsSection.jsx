export const SessionsSection = ({
  sessionStats,
  sessionMessage,
  setSessionMessage,
  onRefresh,
  onDisconnect,
}) => {
  const handleDisconnect = async () => {
    try {
      await onDisconnect()
    } catch (error) {
      setSessionMessage(error.message)
    }
  }

  return (
    <section className="section" id="sessions">
      <div className="section-title">
        <div>
          <h3>Session Control</h3>
          <p>Monitor and disconnect active sessions.</p>
        </div>
        <button className="btn ghost" onClick={onRefresh}>
          Refresh
        </button>
      </div>
      <div className="panel wide">
        <p className="muted small">
          {sessionStats?.session
            ? `Session ${sessionStats.session.id}`
            : 'No active session data yet.'}
        </p>
        <div className="session-grid">
          <div>
            <span className="meta-label">Status</span>
            <span className="meta-value">
              {sessionStats?.session?.status || '—'}
            </span>
          </div>
          <div>
            <span className="meta-label">Expires</span>
            <span className="meta-value">
              {sessionStats?.stats?.expires_at
                ? new Date(sessionStats.stats.expires_at).toLocaleString()
                : '—'}
            </span>
          </div>
          <div>
            <span className="meta-label">Remaining</span>
            <span className="meta-value">
              {sessionStats?.stats?.remaining_minutes
                ? `${sessionStats.stats.remaining_minutes} mins`
                : sessionStats?.stats?.remaining_data_mb
                ? `${sessionStats.stats.remaining_data_mb.toFixed(1)} MB`
                : '—'}
            </span>
          </div>
        </div>
        <button className="btn danger" onClick={handleDisconnect}>
          Disconnect Session
        </button>
        {sessionMessage ? <p className="notice">{sessionMessage}</p> : null}
      </div>
    </section>
  )
}
