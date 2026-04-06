import { useEffect, useState } from 'react';
import { auditAPI } from '../api/client';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ShieldAlert, RefreshCw } from 'lucide-react';

function actionClass(action) {
  if (action.startsWith('phone.'))        return 'bg-blue-50 text-blue-700';
  if (action === 'user.login_failed')     return 'bg-red-50 text-red-700';
  if (action.startsWith('user.'))         return 'bg-green-50 text-green-700';
  if (action.startsWith('cart.'))         return 'bg-yellow-50 text-yellow-700';
  if (action.startsWith('payment.'))      return 'bg-purple-50 text-purple-700';
  return 'bg-muted text-muted-foreground';
}

export default function AuditLog() {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  function load() {
    return auditAPI.list(100)
      .then((data) => setLogs(Array.isArray(data) ? data : []))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }

  useEffect(() => { load(); }, []);

  function handleRefresh() {
    setLoading(true);
    setError('');
    load();
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-8 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
            <ShieldAlert className="h-7 w-7" /> Audit Log
          </h1>
          <p className="text-muted-foreground mt-1">Recent system events across all users and admin actions</p>
        </div>
        <Button variant="outline" size="sm" onClick={handleRefresh} disabled={loading} className="gap-2">
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-base">
            Events <Badge variant="secondary" className="ml-2">{logs.length}</Badge>
          </CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 space-y-2">
              {[...Array(6)].map((_, i) => <div key={i} className="h-8 bg-muted rounded animate-pulse" />)}
            </div>
          ) : error ? (
            <div className="p-6 text-destructive text-sm">{error}</div>
          ) : logs.length === 0 ? (
            <div className="p-6 text-center text-muted-foreground text-sm">No audit events yet.</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-xs">
                <thead>
                  <tr className="border-b bg-muted/40">
                    <th className="px-4 py-2 text-left font-medium text-muted-foreground">Time</th>
                    <th className="px-4 py-2 text-left font-medium text-muted-foreground">Action</th>
                    <th className="px-4 py-2 text-left font-medium text-muted-foreground">User</th>
                    <th className="px-4 py-2 text-left font-medium text-muted-foreground">Resource</th>
                    <th className="px-4 py-2 text-left font-medium text-muted-foreground">IP</th>
                    <th className="px-4 py-2 text-left font-medium text-muted-foreground">Details</th>
                  </tr>
                </thead>
                <tbody className="divide-y">
                  {logs.map((log) => (
                    <tr key={log.id} className="hover:bg-muted/20">
                      <td className="px-4 py-2 whitespace-nowrap text-muted-foreground">
                        {new Date(log.created_at).toLocaleString()}
                      </td>
                      <td className="px-4 py-2 whitespace-nowrap">
                        <span className={`inline-flex items-center rounded px-1.5 py-0.5 font-medium ${actionClass(log.action)}`}>
                          {log.action}
                        </span>
                      </td>
                      <td className="px-4 py-2 text-muted-foreground">{log.user_id ?? '—'}</td>
                      <td className="px-4 py-2 text-muted-foreground">
                        {log.resource_type ? `${log.resource_type}${log.resource_id ? ` #${log.resource_id}` : ''}` : '—'}
                      </td>
                      <td className="px-4 py-2 text-muted-foreground font-mono">{log.ip_address || '—'}</td>
                      <td className="px-4 py-2 text-muted-foreground font-mono max-w-xs truncate">
                        {log.details ? JSON.stringify(log.details) : '—'}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
