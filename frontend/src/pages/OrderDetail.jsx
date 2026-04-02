import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { paymentAPI } from '../api/client';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { ArrowLeft, AlertCircle, Package } from 'lucide-react';

function StatusBadge({ status }) {
  const styles = {
    succeeded: 'bg-green-50 text-green-700 border-green-200',
    pending:   'bg-yellow-50 text-yellow-700 border-yellow-200',
    failed:    'bg-red-50 text-red-700 border-red-200',
  };
  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${styles[status] ?? 'bg-muted text-muted-foreground'}`}>
      {status}
    </span>
  );
}

export default function OrderDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [order, setOrder]   = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError]   = useState('');

  useEffect(() => {
    paymentAPI.getOrder(id)
      .then(setOrder)
      .catch(e => setError(e.message))
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) return (
    <div className="max-w-2xl mx-auto px-4 py-8 space-y-3">
      {[...Array(3)].map((_, i) => (
        <div key={i} className="h-16 bg-muted rounded-lg animate-pulse" />
      ))}
    </div>
  );

  return (
    <div className="max-w-2xl mx-auto px-4 py-8">
      <Button variant="ghost" size="sm" className="mb-6 -ml-2" onClick={() => navigate('/orders')}>
        <ArrowLeft className="h-4 w-4 mr-1" /> Back to Orders
      </Button>

      {error ? (
        <div className="flex items-center gap-2 rounded-md bg-destructive/10 text-destructive px-3 py-2 text-sm">
          <AlertCircle className="h-4 w-4" /> {error}
        </div>
      ) : (
        <>
          <div className="flex items-start justify-between mb-6">
            <div>
              <h1 className="text-2xl font-bold tracking-tight flex items-center gap-2 mb-1">
                <Package className="h-6 w-6" /> Order Details
              </h1>
              <p className="text-sm font-mono text-muted-foreground">{order.id}</p>
            </div>
            <StatusBadge status={order.status} />
          </div>

          {/* Meta */}
          <Card className="mb-4">
            <CardContent className="pt-4 pb-4 grid grid-cols-2 gap-4 text-sm">
              <div>
                <p className="text-xs text-muted-foreground uppercase tracking-wide mb-1">Date</p>
                <p className="font-medium">
                  {order.created_at
                    ? new Date(order.created_at).toLocaleString('en-SG', { dateStyle: 'medium', timeStyle: 'short' })
                    : '—'}
                </p>
              </div>
              <div>
                <p className="text-xs text-muted-foreground uppercase tracking-wide mb-1">Total</p>
                <p className="text-xl font-bold">S${Number(order.total).toFixed(2)}</p>
              </div>
            </CardContent>
          </Card>

          {/* Items */}
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-base font-medium text-muted-foreground">
                {order.items?.length ?? 0} {order.items?.length === 1 ? 'item' : 'items'}
              </CardTitle>
            </CardHeader>
            <CardContent className="p-0">
              {order.items?.map((item, idx) => (
                <div key={item.id ?? idx}>
                  {idx > 0 && <Separator />}
                  <div className="flex items-center justify-between px-6 py-4 gap-4">
                    <div className="flex-1 min-w-0">
                      <p className="font-medium">{item.phone_name || `Phone #${item.phone_id}`}</p>
                      <div className="flex items-center gap-2 mt-1">
                        <Badge variant="secondary" className="text-xs">Qty: {item.quantity}</Badge>
                        <span className="text-sm text-muted-foreground">S${Number(item.price).toFixed(2)} each</span>
                      </div>
                    </div>
                    <span className="font-semibold">
                      S${(Number(item.price) * item.quantity).toFixed(2)}
                    </span>
                  </div>
                </div>
              ))}
              <Separator />
              <div className="flex items-center justify-between px-6 py-4">
                <p className="font-semibold">Total</p>
                <p className="text-2xl font-bold">S${Number(order.total).toFixed(2)}</p>
              </div>
            </CardContent>
          </Card>
        </>
      )}
    </div>
  );
}