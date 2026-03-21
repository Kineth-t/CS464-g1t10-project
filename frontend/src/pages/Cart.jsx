import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { cartAPI } from '../api/client';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { ShoppingCart, Trash2, CheckCircle, AlertCircle, Package } from 'lucide-react';

export default function Cart() {
  const navigate = useNavigate();
  const [cart, setCart] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [feedback, setFeedback] = useState(null);
  const [removingId, setRemovingId] = useState(null);
  const [checkingOut, setCheckingOut] = useState(false);

  function loadCart() {
    return cartAPI.get()
      .then(setCart)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }

  useEffect(() => { loadCart(); }, []);

  async function handleRemove(itemId) {
    setRemovingId(itemId);
    setFeedback(null);
    try {
      await cartAPI.remove(itemId);
      await loadCart();
    } catch (e) {
      setFeedback({ ok: false, msg: e.message });
    } finally {
      setRemovingId(null);
    }
  }

  async function handleCheckout() {
    setCheckingOut(true);
    setFeedback(null);
    try {
      await cartAPI.checkout();
      setFeedback({ ok: true, msg: 'Order placed successfully!' });
      setCart(null);
    } catch (e) {
      setFeedback({ ok: false, msg: e.message });
    } finally {
      setCheckingOut(false);
    }
  }

  if (loading) return (
    <div className="max-w-3xl mx-auto px-4 py-8">
      <div className="animate-pulse space-y-4">
        {[...Array(3)].map((_, i) => <div key={i} className="h-16 bg-muted rounded-lg" />)}
      </div>
    </div>
  );

  return (
    <div className="max-w-3xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold tracking-tight mb-6 flex items-center gap-3">
        <ShoppingCart className="h-7 w-7" /> Your Cart
      </h1>

      {error && (
        <div className="flex items-center gap-2 rounded-md bg-destructive/10 text-destructive px-3 py-2 text-sm mb-4">
          <AlertCircle className="h-4 w-4" /> {error}
        </div>
      )}

      {feedback && (
        <div className={`flex items-center gap-2 rounded-md px-3 py-2 text-sm mb-4 ${feedback.ok ? 'bg-green-50 text-green-700' : 'bg-destructive/10 text-destructive'}`}>
          {feedback.ok ? <CheckCircle className="h-4 w-4" /> : <AlertCircle className="h-4 w-4" />}
          {feedback.msg}
        </div>
      )}

      {(!cart || !cart.items || cart.items.length === 0) ? (
        <div className="flex flex-col items-center gap-3 py-16 text-muted-foreground">
          <Package className="h-12 w-12 opacity-40" />
          <p className="text-lg font-medium">Your cart is empty</p>
          <Button onClick={() => navigate('/')}>Browse Phones</Button>
        </div>
      ) : (
        <Card>
          <CardHeader>
            <CardTitle className="text-base font-medium text-muted-foreground">
              {cart.items.length} {cart.items.length === 1 ? 'item' : 'items'}
            </CardTitle>
          </CardHeader>
          <CardContent className="p-0">
            {cart.items.map((item, idx) => (
              <div key={item.id}>
                {idx > 0 && <Separator />}
                <div className="flex items-center justify-between px-6 py-4 gap-4">
                  <div className="flex-1 min-w-0">
                    <p className="font-medium">Phone #{item.phone_id}</p>
                    <div className="flex items-center gap-2 mt-1">
                      <Badge variant="secondary" className="text-xs">Qty: {item.quantity}</Badge>
                      <span className="text-sm text-muted-foreground">${Number(item.price).toFixed(2)} each</span>
                    </div>
                  </div>
                  <div className="flex items-center gap-4">
                    <span className="font-semibold">${(Number(item.price) * item.quantity).toFixed(2)}</span>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="text-muted-foreground hover:text-destructive"
                      onClick={() => handleRemove(item.id)}
                      disabled={removingId === item.id}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </CardContent>
          <Separator />
          <CardFooter className="flex items-center justify-between px-6 py-4">
            <div>
              <p className="text-sm text-muted-foreground">Total</p>
              <p className="text-2xl font-bold">
                ${cart.items.reduce((sum, item) => sum + Number(item.price) * item.quantity, 0).toFixed(2)}
              </p>
            </div>
            <Button size="lg" onClick={handleCheckout} disabled={checkingOut}>
              {checkingOut ? 'Processing...' : 'Checkout'}
            </Button>
          </CardFooter>
        </Card>
      )}
    </div>
  );
}
