import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { cartAPI, paymentAPI } from '../api/client';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { ArrowLeft, CreditCard, Lock, CheckCircle, AlertCircle, Download } from 'lucide-react';

// ─── Helpers ────────────────────────────────────────────────────────────────

function generateOrderId() {
  return 'ORD-' + Date.now().toString(36).toUpperCase() + '-' +
    Math.random().toString(36).slice(2, 6).toUpperCase();
}

function buildReceiptHTML({ orderId, orderDate, items, total, cardHolder }) {
  const rows = items.map(item => `
    <tr>
      <td style="padding:10px 8px;border-bottom:1px solid #e5e7eb;">Phone #${item.phone_id}</td>
      <td style="padding:10px 8px;border-bottom:1px solid #e5e7eb;text-align:center;">${item.quantity}</td>
      <td style="padding:10px 8px;border-bottom:1px solid #e5e7eb;text-align:right;">
        S$${(Number(item.price) * item.quantity).toFixed(2)}
      </td>
    </tr>
  `).join('');

  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>Receipt ${orderId}</title>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      font-size: 14px;
      color: #111827;
      background: #f9fafb;
      padding: 40px 20px;
    }
    .page {
      max-width: 560px;
      margin: 0 auto;
      background: #ffffff;
      border: 1px solid #e5e7eb;
      border-radius: 12px;
      overflow: hidden;
    }
    .header {
      background: #111827;
      color: #ffffff;
      padding: 28px 32px;
    }
    .header h1 { font-size: 22px; font-weight: 600; margin-bottom: 4px; }
    .header p  { font-size: 13px; color: #9ca3af; }
    .badge {
      display: inline-block;
      background: #d1fae5;
      color: #065f46;
      font-size: 12px;
      font-weight: 600;
      padding: 4px 10px;
      border-radius: 999px;
      margin-top: 12px;
    }
    .body { padding: 28px 32px; }
    .meta { display: flex; justify-content: space-between; margin-bottom: 24px; gap: 16px; flex-wrap: wrap; }
    .meta-block label { font-size: 11px; text-transform: uppercase; letter-spacing: 0.06em; color: #6b7280; display: block; margin-bottom: 4px; }
    .meta-block span  { font-size: 14px; font-weight: 500; }
    table { width: 100%; border-collapse: collapse; margin-bottom: 20px; }
    thead th {
      font-size: 11px;
      text-transform: uppercase;
      letter-spacing: 0.06em;
      color: #6b7280;
      padding: 8px;
      border-bottom: 2px solid #e5e7eb;
      text-align: left;
    }
    thead th:last-child { text-align: right; }
    thead th:nth-child(2) { text-align: center; }
    .total-row td {
      padding: 14px 8px 0;
      font-weight: 600;
      font-size: 16px;
    }
    .total-row td:last-child { text-align: right; }
    .footer {
      border-top: 1px solid #e5e7eb;
      padding: 18px 32px;
      text-align: center;
      font-size: 12px;
      color: #9ca3af;
    }
    @media print {
      body { background: #fff; padding: 0; }
      .page { border: none; border-radius: 0; }
    }
  </style>
</head>
<body>
  <div class="page">
    <div class="header">
      <h1>Ringr Mobile</h1>
      <p>Official purchase receipt</p>
      <span class="badge">Payment confirmed</span>
    </div>
    <div class="body">
      <div class="meta">
        <div class="meta-block">
          <label>Order ID</label>
          <span>${orderId}</span>
        </div>
        <div class="meta-block">
          <label>Date</label>
          <span>${orderDate}</span>
        </div>
        <div class="meta-block">
          <label>Customer</label>
          <span>${cardHolder || '—'}</span>
        </div>
      </div>

      <table>
        <thead>
          <tr>
            <th>Item</th>
            <th>Qty</th>
            <th style="text-align:right;">Amount</th>
          </tr>
        </thead>
        <tbody>
          ${rows}
          <tr class="total-row">
            <td colspan="2">Total</td>
            <td>S$${total.toFixed(2)}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="footer">
      Thank you for shopping with Ringr Mobile &mdash; test mode, no real charges incurred.
    </div>
  </div>
</body>
</html>`;
}

function downloadReceipt({ orderId, orderDate, items, total, cardHolder }) {
  const html = buildReceiptHTML({ orderId, orderDate, items, total, cardHolder });
  const blob = new Blob([html], { type: 'text/html' });
  const url  = URL.createObjectURL(blob);
  const a    = document.createElement('a');
  a.href     = url;
  a.download = `receipt-${orderId}.html`;
  a.click();
  URL.revokeObjectURL(url);
}

// ─── Component ───────────────────────────────────────────────────────────────

export default function Checkout() {
  const navigate = useNavigate();
  const [cart, setCart]               = useState(null);
  const [loadingCart, setLoadingCart] = useState(true);
  const [processing, setProcessing]   = useState(false);
  const [feedback, setFeedback]       = useState(null);
  const [orderMeta, setOrderMeta]     = useState(null); // { orderId, orderDate }

  const [form, setForm] = useState({
    cardHolder: '',
    cardNumber: '',
    expiry: '',
    cvv: '',
  });

  useEffect(() => {
    cartAPI.get()
      .then(setCart)
      .catch(() => navigate('/cart'))
      .finally(() => setLoadingCart(false));
  }, []);

  useEffect(() => {
    if (!loadingCart && (!cart || !cart.items || cart.items.length === 0)) {
      navigate('/cart');
    }
  }, [loadingCart, cart]);

  function handleChange(e) {
    const { name, value } = e.target;
    let formatted = value;

    if (name === 'cardNumber') {
      formatted = value.replace(/\D/g, '').slice(0, 16).replace(/(.{4})/g, '$1 ').trim();
    }
    if (name === 'expiry') {
      formatted = value.replace(/\D/g, '').slice(0, 4);
      if (formatted.length >= 3) formatted = formatted.slice(0, 2) + '/' + formatted.slice(2);
    }
    if (name === 'cvv') {
      formatted = value.replace(/\D/g, '').slice(0, 4);
    }

    setForm(prev => ({ ...prev, [name]: formatted }));
  }

  function getStripeTestMethodId() {
    const raw = form.cardNumber.replace(/\s/g, '');
    const testCards = {
      '4242424242424242': 'pm_card_visa',
      '4000000000000002': 'pm_card_visa_chargeDeclined',
      '4000000000009995': 'pm_card_visa_chargeDeclinedInsufficientFunds',
      '4000000000000069': 'pm_card_visa_chargeDeclinedExpiredCard',
      '4000000000000127': 'pm_card_visa_chargeDeclinedIncorrectCvc',
    };
    return testCards[raw] || 'pm_card_visa';
  }

  async function handlePay(e) {
    e.preventDefault();
    setFeedback(null);

    const raw = form.cardNumber.replace(/\s/g, '');
    if (raw.length !== 16)            return setFeedback({ ok: false, msg: 'Card number must be 16 digits.' });
    if (!form.expiry.match(/^\d{2}\/\d{2}$/)) return setFeedback({ ok: false, msg: 'Expiry must be MM/YY.' });
    if (form.cvv.length < 3)          return setFeedback({ ok: false, msg: 'CVV must be at least 3 digits.' });
    if (!form.cardHolder.trim())      return setFeedback({ ok: false, msg: 'Card holder name is required.' });

    setProcessing(true);
    try {
      const paymentMethodId = getStripeTestMethodId();
      await paymentAPI.pay(paymentMethodId);
      setOrderMeta({
        orderId:   generateOrderId(),
        orderDate: new Date().toLocaleString('en-SG', {
          dateStyle: 'medium',
          timeStyle: 'short',
        }),
      });
    } catch (err) {
      setFeedback({ ok: false, msg: err.message });
    } finally {
      setProcessing(false);
    }
  }

  const total = cart?.items?.reduce(
    (sum, item) => sum + Number(item.price) * item.quantity, 0
  ) ?? 0;

  // ── Success screen ─────────────────────────────────────────────────────────
  if (orderMeta) {
    return (
      <div className="max-w-md mx-auto px-4 py-16 flex flex-col items-center text-center gap-4">
        <div className="rounded-full bg-green-100 p-4">
          <CheckCircle className="h-10 w-10 text-green-600" />
        </div>

        <h1 className="text-2xl font-bold">Payment successful!</h1>

        <p className="text-muted-foreground">
          Your order has been placed. Thank you for shopping with Ringr Mobile.
        </p>

        <div className="bg-muted rounded-lg px-6 py-3 text-sm text-left w-full">
          <p className="text-muted-foreground text-xs uppercase tracking-wide mb-1">Order ID</p>
          <p className="font-mono font-medium">{orderMeta.orderId}</p>
          <p className="text-muted-foreground text-xs mt-1">{orderMeta.orderDate}</p>
        </div>

        <div className="flex flex-col sm:flex-row gap-3 w-full mt-2">
          <Button
            variant="outline"
            className="flex-1 gap-2"
            onClick={() =>
              downloadReceipt({
                orderId:   orderMeta.orderId,
                orderDate: orderMeta.orderDate,
                items:     cart.items,
                total,
                cardHolder: form.cardHolder,
              })
            }
          >
            <Download className="h-4 w-4" />
            Download receipt
          </Button>

          <Button className="flex-1" onClick={() => navigate('/')}>
            Back to Home
          </Button>
        </div>

        <p className="text-xs text-muted-foreground">
          The receipt downloads as an HTML file — open it in any browser and use{' '}
          <kbd className="bg-muted border rounded px-1">Ctrl+P</kbd> to save as PDF.
        </p>
      </div>
    );
  }

  // ── Loading skeleton ───────────────────────────────────────────────────────
  if (loadingCart) return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="animate-pulse space-y-4">
        {[...Array(3)].map((_, i) => <div key={i} className="h-16 bg-muted rounded-lg" />)}
      </div>
    </div>
  );

  // ── Checkout form ──────────────────────────────────────────────────────────
  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <Button variant="ghost" size="sm" className="mb-6 -ml-2" onClick={() => navigate('/cart')}>
        <ArrowLeft className="h-4 w-4 mr-1" /> Back to Cart
      </Button>

      <h1 className="text-3xl font-bold tracking-tight mb-8 flex items-center gap-3">
        <CreditCard className="h-7 w-7" /> Checkout
      </h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">

        {/* Left — card form */}
        <div>
          <Card>
            <CardHeader>
              <CardTitle className="text-base flex items-center gap-2">
                <Lock className="h-4 w-4 text-muted-foreground" /> Payment details
              </CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={handlePay} className="space-y-4">

                <div className="space-y-2">
                  <Label htmlFor="cardHolder">Card holder name</Label>
                  <Input
                    id="cardHolder" name="cardHolder"
                    placeholder="John Doe"
                    value={form.cardHolder} onChange={handleChange} required
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="cardNumber">Card number</Label>
                  <Input
                    id="cardNumber" name="cardNumber"
                    placeholder="1234 5678 9012 3456"
                    value={form.cardNumber} onChange={handleChange}
                    inputMode="numeric" required
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="expiry">Expiry</Label>
                    <Input
                      id="expiry" name="expiry"
                      placeholder="MM/YY"
                      value={form.expiry} onChange={handleChange}
                      inputMode="numeric" required
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="cvv">CVV</Label>
                    <Input
                      id="cvv" name="cvv"
                      placeholder="123"
                      value={form.cvv} onChange={handleChange}
                      inputMode="numeric" type="password" required
                    />
                  </div>
                </div>

                {feedback && (
                  <div className={`flex items-center gap-2 rounded-md px-3 py-2 text-sm ${
                    feedback.ok ? 'bg-green-50 text-green-700' : 'bg-destructive/10 text-destructive'
                  }`}>
                    {feedback.ok
                      ? <CheckCircle className="h-4 w-4" />
                      : <AlertCircle className="h-4 w-4" />}
                    {feedback.msg}
                  </div>
                )}

                <Button type="submit" className="w-full" size="lg" disabled={processing}>
                  <Lock className="h-4 w-4 mr-2" />
                  {processing ? 'Processing payment...' : `Pay S$${total.toFixed(2)}`}
                </Button>

                <p className="text-xs text-muted-foreground text-center">
                  Your payment is secured. Test mode — no real charges.
                </p>
              </form>
            </CardContent>
          </Card>

          {/* Test card hints */}
          <Card className="mt-4 border-dashed">
            <CardContent className="pt-4 pb-4">
              <p className="text-xs font-medium text-muted-foreground mb-2">Test card numbers</p>
              <div className="space-y-1 text-xs text-muted-foreground">
                <p><span className="font-mono">4242 4242 4242 4242</span> — Success</p>
                <p><span className="font-mono">4000 0000 0000 0002</span> — Card declined</p>
                <p><span className="font-mono">4000 0000 0000 9995</span> — Insufficient funds</p>
                <p><span className="font-mono">4000 0000 0000 0069</span> — Expired card</p>
                <p><span className="font-mono">4000 0000 0000 0127</span> — Incorrect CVC</p>
                <p className="mt-2 text-xs">Use any future expiry and any 3-digit CVV.</p>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Right — order summary */}
        <div>
          <Card>
            <CardHeader>
              <CardTitle className="text-base">Order summary</CardTitle>
            </CardHeader>
            <CardContent className="p-0">
              {cart?.items?.map((item, idx) => (
                <div key={item.id}>
                  {idx > 0 && <Separator />}
                  <div className="flex items-center justify-between px-6 py-3 gap-4">
                    <div className="flex-1 min-w-0">
                      <p className="font-medium text-sm">Phone #{item.phone_id}</p>
                      <Badge variant="secondary" className="text-xs mt-1">Qty: {item.quantity}</Badge>
                    </div>
                    <span className="font-semibold text-sm">
                      S${(Number(item.price) * item.quantity).toFixed(2)}
                    </span>
                  </div>
                </div>
              ))}
              <Separator />
              <div className="flex items-center justify-between px-6 py-4">
                <p className="font-semibold">Total</p>
                <p className="text-2xl font-bold">S${total.toFixed(2)}</p>
              </div>
            </CardContent>
          </Card>
        </div>

      </div>
    </div>
  );
}