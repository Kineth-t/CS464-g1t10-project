import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { phonesAPI } from '../api/client';

export default function Home() {
  const [phones, setPhones] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');

  useEffect(() => {
    phonesAPI.list()
      .then(setPhones)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  const filtered = phones.filter((p) => {
    const q = search.toLowerCase();
    return p.brand.toLowerCase().includes(q) || p.model.toLowerCase().includes(q);
  });

  if (loading) return <div className="container"><p className="text-muted">Loading phones...</p></div>;
  if (error) return <div className="container"><p className="text-error">{error}</p></div>;

  return (
    <div className="container">
      <h1 className="page-title">Browse Phones</h1>
      <input
        className="search-input"
        placeholder="Search by brand or model..."
        value={search}
        onChange={(e) => setSearch(e.target.value)}
      />
      {filtered.length === 0 ? (
        <p className="text-muted">No phones found.</p>
      ) : (
        <div className="phone-grid">
          {filtered.map((phone) => (
            <Link to={`/phones/${phone.id}`} key={phone.id} className="phone-card">
              <div className="phone-card-body">
                <h3 className="phone-card-title">{phone.brand} {phone.model}</h3>
                <p className="phone-card-price">${Number(phone.price).toFixed(2)}</p>
                <p className="phone-card-stock">
                  {phone.stock > 0 ? `${phone.stock} in stock` : <span className="out-of-stock">Out of stock</span>}
                </p>
                {phone.description && (
                  <p className="phone-card-desc">{phone.description}</p>
                )}
              </div>
              <div className="phone-card-footer">
                <span className="btn btn-primary btn-sm">View Details</span>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
