import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { cartAPI } from '../api/client';

export default function Cart() {
  const navigate = useNavigate();
  const [cart, setCart] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [feedback, setFeedback] = useState('');

  function loadCart() {
    return cartAPI.get()
      .then(setCart)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }

  useEffect(() => { loadCart(); }, []);

  async function handleRemove(itemId) {
    setFeedback('');
    try {
      await cartAPI.remove(itemId);
      await loadCart();
    } catch (e) {
      setFeedback(e.message);
    }
  }

  async function handleCheckout() {
    setFeedback('');
    try {
      await cartAPI.checkout();
      setFeedback('Order placed successfully!');
      setCart(null);
    } catch (e) {
      setFeedback(e.message);
    }
  }

  if (loading) return <div className="container"><p className="text-muted">Loading cart...</p></div>;
  if (error) return <div className="container"><p className="text-error">{error}</p></div>;

  if (!cart || !cart.items || cart.items.length === 0) {
    return (
      <div className="container">
        <h1 className="page-title">Your Cart</h1>
        {feedback && <p className="text-success">{feedback}</p>}
        <p className="text-muted">Your cart is empty.</p>
        <button className="btn btn-primary" onClick={() => navigate('/')}>Browse Phones</button>
      </div>
    );
  }

  const total = cart.items.reduce((sum, item) => sum + Number(item.price) * item.quantity, 0);

  return (
    <div className="container">
      <h1 className="page-title">Your Cart</h1>
      {feedback && <p className={feedback.includes('!') ? 'text-success' : 'text-error'}>{feedback}</p>}
      <div className="cart-list">
        {cart.items.map((item) => (
          <div key={item.id} className="cart-item">
            <div className="cart-item-info">
              <span className="cart-item-name">Phone #{item.phone_id}</span>
              <span className="cart-item-qty">Qty: {item.quantity}</span>
              <span className="cart-item-price">${(Number(item.price) * item.quantity).toFixed(2)}</span>
            </div>
            <button className="btn btn-danger btn-sm" onClick={() => handleRemove(item.id)}>Remove</button>
          </div>
        ))}
      </div>
      <div className="cart-footer">
        <span className="cart-total">Total: ${total.toFixed(2)}</span>
        <button className="btn btn-primary" onClick={handleCheckout}>Checkout</button>
      </div>
    </div>
  );
}
