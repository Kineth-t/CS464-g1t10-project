import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { cartAPI, paymentAPI } from '../api/client';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { ArrowLeft, CreditCard, Lock, CheckCircle, AlertCircle } from 'lucide-react';

export default function Checkout() {
  const navigate = useNavigate();
  const [cart, setCart] = useState(null);
  const [loadingCart, setLoadingCart] = useState(true);
  const [processing, setProcessing] = useState(false);
  const [feedback, setFeedback] = useState(null);
  const [success, setSuccess] = useState(false);

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

  // Redirect back to cart if cart is empty
  useEffect(() => {
    if (!loadingCart && (!cart || !cart.items || cart.items.length === 0)) {
      navigate('/cart');
    }
  }, [loadingCart, cart]);

  function handleChange(e) {
    const { name, value } = e.target;
    let formatted = value;

    // Format card number as XXXX XXXX XXXX XXXX
    if (name === 'cardNumber') {
      formatted = value.replace(/\D/g, '').slice(0, 16).replace(/(.{4})/g, '$1 ').trim();
    }

    // Format expiry as MM/YY
    if (name === 'expiry') {
      formatted = value.replace(/\D/g, '').slice(0, 4);
      if (formatted.length >= 3) {
        formatted = formatted.slice(0, 2) + '/' + formatted.slice(2);
      }
    }

    // CVV max 4 digits
    if (name === 'cvv') {
      formatted = value.replace(/\D/g, '').slice(0, 4);
    }

    setForm((prev) => ({ ...prev, [name]: formatted }));
  }

  function getStripeTestMethodId() {
    // Map the entered card number to a Stripe test payment method ID
    // In production the frontend would use Stripe.js to tokenise the real card
    const raw = form.cardNumber.replace(/\s/g, '');

    const testCards = {
      '4242424242424242': 'pm_card_visa',
      '4000000000000002': 'pm_card_visa_chargeDeclined',
      '4000000000009995': 'pm_card_visa_chargeDeclinedInsufficientFunds',
      '4000000000000069': 'pm_card_visa_chargeDeclinedExpiredCard',
      '4000000000000127': 'pm_card_visa_chargeDeclinedIncorrectCvc',
    };

    return testCards[raw] || 'pm_card_visa'; // default to success card
  }

  async function handlePay(e) {
    e.preventDefault();
    setFeedback(null);

    // Basic validation
    const raw = form.cardNumber.replace(/\s/g, '');
    if (raw.length !== 16) return setFeedback({ ok: false, msg: 'Card number must be 16 digits.' });
    if (!form.expiry.match(/^\d{2}\/\d{2}$/)) return setFeedback({ ok: false, msg: 'Expiry must be MM/YY.' });
    if (form.cvv.length < 3) return setFeedback({ ok: false, msg: 'CVV must be at least 3 digits.' });
    if (!form.cardHolder.trim()) return setFeedback({ ok: false, msg: 'Card holder name is required.' });

    setProcessing(true);
    try {
      // Get the Stripe test payment method ID based on the test card entered
      const paymentMethodId = getStripeTestMethodId();
      await paymentAPI.pay(paymentMethodId);
      setSuccess(true);
    } catch (e) {
      setFeedback({ ok: false, msg: e.message });
    } finally {
      setProcessing(false);
    }
  }

  const total = cart?.items?.reduce(
    (sum, item) => sum + Number(item.price) * item.quantity, 0
  ) ?? 0;

  // Success screen
  if (success) {
    return (
      <div className="max-w-md mx-auto px-4 py-16 flex flex-col items-center text-center gap-4">
        <div className="rounded-full bg-green-100 p-4">
          <CheckCircle className="h-10 w-10 text-green-600" />
        </div>
        <h1 className="text-2xl font-bold">Payment successful!</h1>
        <p className="text-muted-foreground">Your order has been placed. Thank you for shopping with Ringr Mobile.</p>
        <Button className="mt-4" onClick={() => navigate('/')}>Back to Home</Button>
      </div>
    );
  }

  if (loadingCart) return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="animate-pulse space-y-4">
        {[...Array(3)].map((_, i) => <div key={i} className="h-16 bg-muted rounded-lg" />)}
      </div>
    </div>
  );

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
                    id="cardHolder"
                    name="cardHolder"
                    placeholder="John Doe"
                    value={form.cardHolder}
                    onChange={handleChange}
                    required
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="cardNumber">Card number</Label>
                  <Input
                    id="cardNumber"
                    name="cardNumber"
                    placeholder="1234 5678 9012 3456"
                    value={form.cardNumber}
                    onChange={handleChange}
                    inputMode="numeric"
                    required
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="expiry">Expiry</Label>
                    <Input
                      id="expiry"
                      name="expiry"
                      placeholder="MM/YY"
                      value={form.expiry}
                      onChange={handleChange}
                      inputMode="numeric"
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="cvv">CVV</Label>
                    <Input
                      id="cvv"
                      name="cvv"
                      placeholder="123"
                      value={form.cvv}
                      onChange={handleChange}
                      inputMode="numeric"
                      type="password"
                      required
                    />
                  </div>
                </div>

                {feedback && (
                  <div className={`flex items-center gap-2 rounded-md px-3 py-2 text-sm ${feedback.ok ? 'bg-green-50 text-green-700' : 'bg-destructive/10 text-destructive'}`}>
                    {feedback.ok ? <CheckCircle className="h-4 w-4" /> : <AlertCircle className="h-4 w-4" />}
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