import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function Navbar() {
  const { isAuthenticated, isAdmin, user, logout } = useAuth();
  const navigate = useNavigate();

  function handleLogout() {
    logout();
    navigate('/login');
  }

  return (
    <nav className="navbar">
      <Link to="/" className="navbar-brand">PhoneStore</Link>
      <div className="navbar-links">
        <Link to="/">Phones</Link>
        {isAuthenticated && <Link to="/cart">Cart</Link>}
        {isAdmin && <Link to="/admin">Admin</Link>}
        {isAuthenticated ? (
          <span className="navbar-user">
            {user?.username}
            <button onClick={handleLogout} className="btn btn-outline btn-sm">Logout</button>
          </span>
        ) : (
          <>
            <Link to="/login">Login</Link>
            <Link to="/register" className="btn btn-primary btn-sm">Register</Link>
          </>
        )}
      </div>
    </nav>
  );
}
