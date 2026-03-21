import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const INITIAL = {
  username: '', password: '', phone_number: '',
  street: '', city: '', state: '', country: '', zip_code: '',
};

export default function Register() {
  const { register } = useAuth();
  const navigate = useNavigate();
  const [form, setForm] = useState(INITIAL);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  function handleChange(e) {
    setForm((f) => ({ ...f, [e.target.name]: e.target.value }));
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await register(form);
      navigate('/login');
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  function field(label, name, type = 'text', required = false) {
    return (
      <label className="form-label" key={name}>
        {label}{required && ' *'}
        <input className="form-input" type={type} name={name} value={form[name]} onChange={handleChange} required={required} />
      </label>
    );
  }

  return (
    <div className="auth-page">
      <form className="auth-form auth-form--wide" onSubmit={handleSubmit}>
        <h2 className="auth-title">Create Account</h2>
        {error && <p className="text-error">{error}</p>}
        <fieldset className="form-section">
          <legend>Account</legend>
          {field('Username', 'username', 'text', true)}
          {field('Password', 'password', 'password', true)}
          {field('Phone Number', 'phone_number')}
        </fieldset>
        <fieldset className="form-section">
          <legend>Address</legend>
          {field('Street', 'street')}
          {field('City', 'city')}
          {field('State', 'state')}
          {field('Country', 'country')}
          {field('Zip Code', 'zip_code')}
        </fieldset>
        <button className="btn btn-primary btn-full" type="submit" disabled={loading}>
          {loading ? 'Creating account...' : 'Register'}
        </button>
        <p className="auth-switch">Already have an account? <Link to="/login">Login</Link></p>
      </form>
    </div>
  );
}
