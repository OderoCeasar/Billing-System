import { useState } from 'react'

export const usePayments = (apiFetch) => {
  const [paymentForm, setPaymentForm] = useState({
    package_id: '',
    phone_number: '',
  })
  const [paymentMessage, setPaymentMessage] = useState('')

  const initiatePayment = async () => {
    setPaymentMessage('')
    const data = await apiFetch('/payments/initiate', {
      method: 'POST',
      body: JSON.stringify(paymentForm),
    })
    const paymentId = data?.payment_id || data?.payment?.id || 'payment'
    setPaymentMessage(`STK push sent. Payment ID: ${paymentId}`)
    return data
  }

  return {
    paymentForm,
    setPaymentForm,
    paymentMessage,
    setPaymentMessage,
    initiatePayment,
  }
}
