import { useEffect, useState } from 'react';
import { phonesAPI } from '../api/client';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { PlusCircle, Pencil, Trash2, CheckCircle, AlertCircle, X } from 'lucide-react';

const EMPTY_FORM = { brand: '', model: '', price: '', stock: '', description: '' };

export default function Admin() {
  const [phones, setPhones] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [feedback, setFeedback] = useState(null);
  const [form, setForm] = useState(EMPTY_FORM);
  const [editingId, setEditingId] = useState(null);
  const [submitting, setSubmitting] = useState(false);

  function loadPhones() {
    return phonesAPI.list()
      .then(setPhones)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }

  useEffect(() => { loadPhones(); }, []);

  function handleChange(e) {
    setForm((f) => ({ ...f, [e.target.name]: e.target.value }));
  }

  function startEdit(phone) {
    setEditingId(phone.id);
    setForm({ brand: phone.brand, model: phone.model, price: String(phone.price), stock: String(phone.stock), description: phone.description || '' });
    setFeedback(null);
  }

  function cancelEdit() {
    setEditingId(null);
    setForm(EMPTY_FORM);
    setFeedback(null);
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setFeedback(null);
    setSubmitting(true);
    const payload = { brand: form.brand, model: form.model, price: parseFloat(form.price), stock: parseInt(form.stock, 10), description: form.description };
    try {
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
    </div>
  );
}
