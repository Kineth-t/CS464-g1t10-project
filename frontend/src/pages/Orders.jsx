import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { paymentAPI } from '../api/client';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Package, ChevronRight, ShoppingBag, AlertCircle } from 'lucide-react';

function StatusBadge({ status }) {
  const styles = {
    succeeded: 'bg-green-50 text-green-700 border-green-200',
    pending:   'bg-yellow-50 text-yellow-700 border-yellow-200',
    failed:    'bg-red-50 text-red-700 border-red-200',
  };
  return (
    <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium border ${styles[status] ?? 'bg-muted text-muted-foreground'}`}>
      {status}
    </span>
  );
}

export default function Orders() {
  const navigate = useNavigate();
  const [orders, setOrders]   = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError]     = useState('');

  useEffect(() => {
    paymentAPI.getOrders()
      .then(data => setOrders(data ?? []))
      .catch(e => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return (
    <div className="max-w-3xl mx-auto px-4 py-8 space-y-3">
      {[...Array(3)].map((_, i) => (
        <div key={i} className="h-20 bg-muted rounded-lg animate-pulse" />
      ))}
    </div>
  );

  return (
    <div className="max-w-3xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold tracking-tight mb-6 flex items-center gap-3">
        <ShoppingBag className="h-7 w-7" /> Your Orders
      </h1>

      {error && (
        <div className="flex items-center gap-2 rounded-md bg-destructive/10 text-destructive px-3 py-2 text-sm mb-4">
          <AlertCircle className="h-4 w-4" /> {error}
        </div>
      )}

      {orders.length === 0 ? (
        <div className="flex flex-col items-center gap-3 py-16 text-muted-foreground">
          <Package className="h-12 w-12 opacity-40" />
          <p className="text-lg font-medium">No orders yet</p>
          <p className="text-sm">Your completed orders will appear here.</p>
        </div>
      ) : (
        <Card>
          {orders.map((order, idx) => (
            <div key={order.id}>
              {idx > 0 && <Separator />}
              <button
                className="w-full text-left px-6 py-4 flex items-center justify-between gap-4 hover:bg-muted/50 transition-colors"
                onClick={() => navigate(`/orders/${order.id}`)}
              >
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-1">
                    <StatusBadge status={order.status} />
                    <span className="text-xs text-muted-foreground font-mono truncate">{order.id}</span>
                  </div>
                  <p className="text-sm text-muted-foreground">
                    {order.items?.length ?? 0} {order.items?.length === 1 ? 'item' : 'items'}
                    {order.created_at && (
                      <> &middot; {new Date(order.created_at).toLocaleDateString('en-SG', { dateStyle: 'medium' })}</>
                    )}
                  </p>
                </div>
                <div className="flex items-center gap-3">
                  <span className="font-semibold">S${Number(order.total).toFixed(2)}</span>
                  <ChevronRight className="h-4 w-4 text-muted-foreground" />
                </div>
              </button>
            </div>
          ))}
        </Card>
      )}
    </div>
  );
}