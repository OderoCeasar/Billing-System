export const PaymentsSection = ({
  packages,
  token,
  paymentForm,
  setPaymentForm,
  paymentMessage,
  setPaymentMessage,
  onSubmit,
}) => {
  const handleSubmit = async (event) => {
    event.preventDefault()
    try {
      await onSubmit()
    } catch (error) {
      setPaymentMessage(error.message)
    }
  }

  return (
    <section className="section" id="payments">
      <div className="section-title">
        <div>
          <h3>Payments</h3>
          <p>Send STK push for a selected package.</p>
        </div>
      </div>
      <form className="panel wide" onSubmit={handleSubmit}>
        <label>
          Package
          <select
            value={paymentForm.package_id}
            onChange={(e) =>
              setPaymentForm((prev) => ({
                ...prev,
                package_id: e.target.value,
              }))
            }
          >
            <option value="">Select a package</option>
            {packages.map((pkg) => (
              <option key={pkg.id} value={pkg.id}>
                {pkg.name} - Ksh {pkg.price}
              </option>
            ))}
          </select>
        </label>
        <label>
          Phone Number
          <input
            type="tel"
            value={paymentForm.phone_number}
            onChange={(e) =>
              setPaymentForm((prev) => ({
                ...prev,
                phone_number: e.target.value,
              }))
            }
            placeholder="07XXXXXXXX"
          />
        </label>
        <button className="btn primary" type="submit" disabled={!token}>
          Initiate Payment
        </button>
        {!token ? (
          <p className="muted small">Login required to initiate payments.</p>
        ) : null}
      </form>
      {paymentMessage ? <div className="notice">{paymentMessage}</div> : null}
    </section>
  )
}
