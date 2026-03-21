import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { phonesAPI, cartAPI } from '../api/client';
import { useAuth } from '../context/AuthContext';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { ArrowLeft, ShoppingCart, CheckCircle, AlertCircle } from 'lucide-react';

export default function PhoneDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const [phone, setPhone] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [quantity, setQuantity] = useState(1);
  const [feedback, setFeedback] = useState(null);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    phonesAPI.get(id)
      .then(setPhone)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [id]);

  async function handleAddToCart() {
    if (!isAuthenticated) { navigate('/login'); return; }
    setSubmitting(true);
    setFeedback(null);
    try {
      await cartAPI.add(phone.id, quantity);
      setFeedback({ ok: true, msg: 'Added to cart!' });
    } catch (e) {
      setFeedback({ ok: false, msg: e.message });
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <div className="animate-pulse space-y-4">
        <div className="h-6 bg-muted rounded w-24" />
        <div className="h-8 bg-muted rounded w-64" />
        <div className="h-4 bg-muted rounded w-32" />
      </div>
    </div>
  );

  if (error) return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-4 text-destructive text-sm flex items-center gap-2">
        <AlertCircle className="h-4 w-4" /> {error}
      </div>
    </div>
  );

  if (!phone) return null;
  const outOfStock = phone.stock === 0;

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <Button variant="ghost" size="sm" className="mb-6 -ml-2" onClick={() => navigate(-1)}>
        <ArrowLeft className="h-4 w-4 mr-1" /> Back
      </Button>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div className="md:col-span-2 space-y-4">
          <div>
            <p className="text-sm text-muted-foreground font-medium uppercase tracking-wider">{phone.brand}</p>
            <h1 className="text-3xl font-bold mt-1">{phone.model}</h1>
          </div>
          <div className="flex items-center gap-3">
            <span className="text-4xl font-bold">${Number(phone.price).toFixed(2)}</span>
            {outOfStock
              ? <Badge variant="destructive">Out of stock</Badge>
              : <Badge variant="secondary">{phone.stock} units available</Badge>
            }
          </div>
          {phone.description && (
            <>
              <Separator />
              <p className="text-muted-foreground leading-relaxed">{phone.description}</p>
            </>
          )}
        </div>

        <div>
          <Card>
            <CardContent className="pt-6 space-y-4">
              <div className="space-y-2">
                <Label htmlFor="qty">Quantity</Label>
                <Input
                  id="qty"
                  type="number"
                  min={1}
                  max={phone.stock}
                  value={quantity}
                  disabled={outOfStock}
                  onChange={(e) => setQuantity(Math.max(1, Number(e.target.value)))}
                  className="w-24"
                />
              </div>
              <Button
                className="w-full"
                onClick={handleAddToCart}
                disabled={outOfStock || submitting}
              >
                <ShoppingCart className="h-4 w-4 mr-2" />
                {submitting ? 'Adding...' : outOfStock ? 'Out of Stock' : 'Add to Cart'}
              </Button>
              {feedback && (
                <div className={`flex items-center gap-2 text-sm ${feedback.ok ? 'text-green-600' : 'text-destructive'}`}>
                  {feedback.ok ? <CheckCircle className="h-4 w-4" /> : <AlertCircle className="h-4 w-4" />}
                  {feedback.msg}
                </div>
              )}
              {!isAuthenticated && (
                <p className="text-xs text-muted-foreground text-center">
                  <button onClick={() => navigate('/login')} className="underline">Login</button> to add to cart
                </p>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
