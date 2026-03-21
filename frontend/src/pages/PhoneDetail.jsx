import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { phonesAPI, cartAPI } from '../api/client';
import { useAuth } from '../context/AuthContext';

export default function PhoneDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const [phone, setPhone] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [quantity, setQuantity] = useState(1);
  const [feedback, setFeedback] = useState('');
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
    setFeedback('');
    try {
      await cartAPI.add(phone.id, quantity);
      setFeedback('Added to cart!');
    } catch (e) {
      setFeedback(e.message);
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) return <div className="container"><p className="text-muted">Loading...</p></div>;
  if (error) return <div className="container"><p className="text-error">{error}</p></div>;
  if (!phone) return null;

  const outOfStock = phone.stock === 0;

  return (
    <div className="container">
      <button className="btn btn-outline btn-sm back-btn" onClick={() => navigate(-1)}>← Back</button>
      <div className="phone-detail">
        <div className="phone-detail-info">
          <h1 className="phone-detail-title">{phone.brand} {phone.model}</h1>
          <p className="phone-detail-price">${Number(phone.price).toFixed(2)}</p>
          <p className={`phone-detail-stock ${outOfStock ? 'out-of-stock' : 'in-stock'}`}>
            {outOfStock ? 'Out of stock' : `${phone.stock} units available`}
          </p>
          {phone.description && <p className="phone-detail-desc">{phone.description}</p>}
        </div>
        <div className="phone-detail-actions">
          <label className="qty-label">
            Quantity
            <input
              type="number"
              className="qty-input"
              min={1}
              max={phone.stock}
              value={quantity}
              disabled={outOfStock}
              onChange={(e) => setQuantity(Math.max(1, Number(e.target.value)))}
            />
          </label>
          <button
            className="btn btn-primary"
            onClick={handleAddToCart}
            disabled={outOfStock || submitting}
          >
            {submitting ? 'Adding...' : 'Add to Cart'}
          </button>
          {feedback && <p className={feedback.includes('!') ? 'text-success' : 'text-error'}>{feedback}</p>}
        </div>
      </div>
    </div>
  );
}
