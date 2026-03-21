import { useEffect, useState } from 'react';
import { phonesAPI } from '../api/client';

const EMPTY_FORM = { brand: '', model: '', price: '', stock: '', description: '' };

export default function Admin() {
  const [phones, setPhones] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [feedback, setFeedback] = useState('');
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
    setForm({
      brand: phone.brand,
      model: phone.model,
      price: String(phone.price),
      stock: String(phone.stock),
      description: phone.description || '',
    });
    setFeedback('');
  }

  function cancelEdit() {
    setEditingId(null);
    setForm(EMPTY_FORM);
    setFeedback('');
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setFeedback('');
    setSubmitting(true);
    const payload = {
      brand: form.brand,
      model: form.model,
      price: parseFloat(form.price),
      stock: parseInt(form.stock, 10),
      description: form.description,
    };
    try {
      if (editingId) {
        await phonesAPI.update(editingId, payload);
        setFeedback('Phone updated.');
      } else {
        await phonesAPI.create(payload);
        setFeedback('Phone created.');
      }
      cancelEdit();
      await loadPhones();
    } catch (err) {
      setFeedback(err.message);
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(id) {
    if (!confirm('Delete this phone?')) return;
    setFeedback('');
    try {
      await phonesAPI.delete(id);
      setFeedback('Phone deleted.');
      await loadPhones();
    } catch (err) {
      setFeedback(err.message);
    }
  }

  return (
    <div className="container">
      <h1 className="page-title">Admin — Phone Management</h1>
      {feedback && (
        <p className={feedback.toLowerCase().includes('error') || feedback.toLowerCase().includes('fail') ? 'text-error' : 'text-success'}>
          {feedback}
        </p>
      )}

      <form className="admin-form" onSubmit={handleSubmit}>
        <h2 className="section-title">{editingId ? 'Edit Phone' : 'Add New Phone'}</h2>
        <div className="admin-form-grid">
          {[
            ['Brand', 'brand', 'text'],
            ['Model', 'model', 'text'],
            ['Price', 'price', 'number'],
            ['Stock', 'stock', 'number'],
          ].map(([label, name, type]) => (
            <label className="form-label" key={name}>
              {label}
              <input
                className="form-input"
                type={type}
                name={name}
                value={form[name]}
                onChange={handleChange}
                min={type === 'number' ? 0 : undefined}
                step={name === 'price' ? '0.01' : undefined}
                required
              />
            </label>
          ))}
        </div>
        <label className="form-label">
          Description
          <textarea className="form-input form-textarea" name="description" value={form.description} onChange={handleChange} />
        </label>
        <div className="admin-form-actions">
          <button className="btn btn-primary" type="submit" disabled={submitting}>
            {submitting ? 'Saving...' : editingId ? 'Update Phone' : 'Create Phone'}
          </button>
          {editingId && (
            <button className="btn btn-outline" type="button" onClick={cancelEdit}>Cancel</button>
          )}
        </div>
      </form>

      {loading ? (
        <p className="text-muted">Loading phones...</p>
      ) : error ? (
        <p className="text-error">{error}</p>
      ) : (
        <table className="admin-table">
          <thead>
            <tr>
              <th>ID</th><th>Brand</th><th>Model</th><th>Price</th><th>Stock</th><th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {phones.map((phone) => (
              <tr key={phone.id} className={editingId === phone.id ? 'row-editing' : ''}>
                <td>{phone.id}</td>
                <td>{phone.brand}</td>
                <td>{phone.model}</td>
                <td>${Number(phone.price).toFixed(2)}</td>
                <td>{phone.stock}</td>
                <td className="admin-actions">
                  <button className="btn btn-outline btn-sm" onClick={() => startEdit(phone)}>Edit</button>
                  <button className="btn btn-danger btn-sm" onClick={() => handleDelete(phone.id)}>Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
