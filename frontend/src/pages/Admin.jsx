import { useEffect, useRef, useState } from 'react';
import { auditAPI, phonesAPI, uploadAPI } from '../api/client';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { PlusCircle, Pencil, Trash2, CheckCircle, AlertCircle, X, RefreshCw, ShieldAlert } from 'lucide-react';

const EMPTY_FORM = { brand: '', model: '', price: '', stock: '', description: '', image_url: '' };

export default function Admin() {
  const [phones, setPhones] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [feedback, setFeedback] = useState(null);
  const [form, setForm] = useState(EMPTY_FORM);
  const [editingId, setEditingId] = useState(null);
  const [submitting, setSubmitting] = useState(false);
  const [auditLogs, setAuditLogs] = useState([]);
  const [auditLoading, setAuditLoading] = useState(true);
  const [auditError, setAuditError] = useState('');
  const [imageFile, setImageFile] = useState(null);
  const [imagePreview, setImagePreview] = useState('');
  const fileInputRef = useRef(null);

  function loadPhones() {
    return phonesAPI.list()
      .then(setPhones)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }

  function loadAuditLogs() {
    setAuditLoading(true);
    setAuditError('');
    return auditAPI.list(100)
      .then((data) => setAuditLogs(Array.isArray(data) ? data : []))
      .catch((e) => setAuditError(e.message))
      .finally(() => setAuditLoading(false));
  }

  useEffect(() => { loadPhones(); loadAuditLogs(); }, []);

  function handleChange(e) {
    setForm((f) => ({ ...f, [e.target.name]: e.target.value }));
  }

  function handleImageChange(e) {
    const file = e.target.files[0];
    if (!file) return;
    setImageFile(file);
    setImagePreview(URL.createObjectURL(file));
  }

  function clearImage() {
    setImageFile(null);
    setImagePreview('');
    setForm((f) => ({ ...f, image_url: '' }));
    if (fileInputRef.current) fileInputRef.current.value = '';
  }

  function startEdit(phone) {
    setEditingId(phone.id);
    setForm({ brand: phone.brand, model: phone.model, price: String(phone.price), stock: String(phone.stock), description: phone.description || '', image_url: phone.image_url || '' });
    setImageFile(null);
    setImagePreview(phone.image_url || '');
    if (fileInputRef.current) fileInputRef.current.value = '';
    setFeedback(null);
  }

  function cancelEdit() {
    setEditingId(null);
    setForm(EMPTY_FORM);
    setImageFile(null);
    setImagePreview('');
    if (fileInputRef.current) fileInputRef.current.value = '';
    setFeedback(null);
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setFeedback(null);
    setSubmitting(true);
    try {
      let imageUrl = form.image_url;

      // If a new file was selected, upload it first
      if (imageFile) {
        const result = await uploadAPI.upload(imageFile);
        imageUrl = result.url;
      }

      const payload = {
        brand: form.brand,
        model: form.model,
        price: parseFloat(form.price),
        stock: parseInt(form.stock, 10),
        description: form.description,
        image_url: imageUrl,
      };

      if (editingId) {
        await phonesAPI.update(editingId, payload);
        setFeedback({ ok: true, msg: 'Phone updated.' });
      } else {
        await phonesAPI.create(payload);
        setFeedback({ ok: true, msg: 'Phone created.' });
      }
      cancelEdit();
      await loadPhones();
    } catch (err) {
      setFeedback({ ok: false, msg: err.message });
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(id) {
    if (!confirm('Delete this phone?')) return;
    setFeedback(null);
    try {
      await phonesAPI.delete(id);
      setFeedback({ ok: true, msg: 'Phone deleted.' });
      await loadPhones();
    } catch (err) {
      setFeedback({ ok: false, msg: err.message });
    }
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-8 space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Phone Management</h1>
        <p className="text-muted-foreground mt-1">Add, edit, or remove phones from the catalog</p>
      </div>

      {feedback && (
        <div className={`flex items-center gap-2 rounded-md px-3 py-2 text-sm ${feedback.ok ? 'bg-green-50 text-green-700' : 'bg-destructive/10 text-destructive'}`}>
          {feedback.ok ? <CheckCircle className="h-4 w-4" /> : <AlertCircle className="h-4 w-4" />}
          {feedback.msg}
        </div>
      )}

      <Card>
        <CardHeader className="pb-4">
          <CardTitle className="flex items-center gap-2 text-base">
            <PlusCircle className="h-4 w-4" />
            {editingId ? 'Edit Phone' : 'Add New Phone'}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {[['Brand', 'brand', 'text'], ['Model', 'model', 'text'], ['Price ($)', 'price', 'number'], ['Stock', 'stock', 'number']].map(([label, name, type]) => (
                <div key={name} className="space-y-2">
                  <Label htmlFor={name}>{label}</Label>
                  <Input id={name} name={name} type={type} value={form[name]} onChange={handleChange}
                    min={type === 'number' ? 0 : undefined} step={name === 'price' ? '0.01' : undefined} required />
                </div>
              ))}
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Description</Label>
              <Input id="description" name="description" value={form.description} onChange={handleChange} placeholder="Optional" />
            </div>
            <div className="space-y-2">
              <Label>Product Image</Label>
              <div className="flex items-start gap-4">
                {imagePreview && (
                  <div className="relative shrink-0">
                    <img src={imagePreview} alt="Preview" className="h-20 w-20 rounded-md object-cover border" />
                    <button type="button" onClick={clearImage} className="absolute -top-1.5 -right-1.5 bg-destructive text-destructive-foreground rounded-full h-5 w-5 flex items-center justify-center text-xs leading-none">
                      <X className="h-3 w-3" />
                    </button>
                  </div>
                )}
                <div className="flex-1 space-y-1.5">
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/*"
                    onChange={handleImageChange}
                    className="block w-full text-sm text-muted-foreground file:mr-3 file:py-1.5 file:px-3 file:rounded-md file:border-0 file:text-sm file:font-medium file:bg-secondary file:text-secondary-foreground hover:file:bg-secondary/80 cursor-pointer"
                  />
                  <p className="text-xs text-muted-foreground">Upload a new image, or paste a URL below</p>
                  <Input
                    id="image_url"
                    name="image_url"
                    value={imageFile ? '' : form.image_url}
                    onChange={(e) => { setImageFile(null); setImagePreview(e.target.value); setForm((f) => ({ ...f, image_url: e.target.value })); }}
                    placeholder="https://example.com/image.jpg"
                    disabled={!!imageFile}
                  />
                </div>
              </div>
            </div>
            <div className="flex gap-2">
              <Button type="submit" disabled={submitting}>
                {submitting ? 'Saving...' : editingId ? 'Update Phone' : 'Create Phone'}
              </Button>
              {editingId && (
                <Button type="button" variant="outline" onClick={cancelEdit}>
                  <X className="h-4 w-4 mr-1" /> Cancel
                </Button>
              )}
            </div>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-base">All Phones <Badge variant="secondary" className="ml-2">{phones.length}</Badge></CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-6 space-y-3">
              {[...Array(3)].map((_, i) => <div key={i} className="h-10 bg-muted rounded animate-pulse" />)}
            </div>
          ) : error ? (
            <div className="p-6 text-destructive text-sm">{error}</div>
          ) : phones.length === 0 ? (
            <div className="p-6 text-center text-muted-foreground text-sm">No phones yet. Add one above.</div>
          ) : (
            <div className="divide-y">
              {phones.map((phone) => (
                <div key={phone.id} className={`flex items-center justify-between px-6 py-3 gap-4 ${editingId === phone.id ? 'bg-blue-50' : 'hover:bg-muted/40'}`}>
                  <div className="flex items-center gap-4 min-w-0 flex-1">
                    <span className="text-xs text-muted-foreground w-6">#{phone.id}</span>
                    <span className="font-medium truncate">{phone.brand} {phone.model}</span>
                  </div>
                  <div className="flex items-center gap-4 shrink-0">
                    <span className="text-sm font-semibold">${Number(phone.price).toFixed(2)}</span>
                    <Badge variant={phone.stock > 0 ? 'secondary' : 'destructive'}>{phone.stock} in stock</Badge>
                    <div className="flex gap-1">
                      <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => startEdit(phone)}>
                        <Pencil className="h-3.5 w-3.5" />
                      </Button>
                      <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-destructive" onClick={() => handleDelete(phone.id)}>
                        <Trash2 className="h-3.5 w-3.5" />
                      </Button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="pb-2">
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2 text-base">
              <ShieldAlert className="h-4 w-4" />
              Audit Log
              <Badge variant="secondary" className="ml-1">{auditLogs.length}</Badge>
            </CardTitle>
            <Button variant="ghost" size="icon" className="h-8 w-8" onClick={loadAuditLogs} disabled={auditLoading}>
              <RefreshCw className={`h-3.5 w-3.5 ${auditLoading ? 'animate-spin' : ''}`} />
            </Button>
          </div>
        </CardHeader>
        <CardContent className="p-0">
          {auditLoading ? (
            <div className="p-6 space-y-2">
              {[...Array(4)].map((_, i) => <div key={i} className="h-8 bg-muted rounded animate-pulse" />)}
            </div>
          ) : auditError ? (
            <div className="p-6 text-destructive text-sm">{auditError}</div>
          ) : auditLogs.length === 0 ? (
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
                  {auditLogs.map((log) => (
                    <tr key={log.id} className="hover:bg-muted/20">
                      <td className="px-4 py-2 whitespace-nowrap text-muted-foreground">
                        {new Date(log.created_at).toLocaleString()}
                      </td>
                      <td className="px-4 py-2 whitespace-nowrap">
                        <span className={`inline-flex items-center rounded px-1.5 py-0.5 font-medium ${
                          log.action.startsWith('phone.') ? 'bg-blue-50 text-blue-700' :
                          log.action === 'user.login_failed' ? 'bg-red-50 text-red-700' :
                          log.action.startsWith('user.') ? 'bg-green-50 text-green-700' :
                          log.action.startsWith('cart.') ? 'bg-yellow-50 text-yellow-700' :
                          log.action.startsWith('payment.') ? 'bg-purple-50 text-purple-700' :
                          'bg-muted text-muted-foreground'
                        }`}>
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
