import { createContext, useContext, useState, useCallback } from 'react';
import { authAPI } from '../api/client';

const AuthContext = createContext(null);

function parseJwt(token) {
  try {
    return JSON.parse(atob(token.split('.')[1]));
  } catch {
    return null;
  }
}

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => localStorage.getItem('token'));
  const [user, setUser] = useState(() => {
    const t = localStorage.getItem('token');
    return t ? parseJwt(t) : null;
  });

  const login = useCallback(async (username, password) => {
    const data = await authAPI.login(username, password);
    localStorage.setItem('token', data.token);
    setToken(data.token);
    setUser(parseJwt(data.token));
  }, []);

  const register = useCallback(async (payload) => {
    await authAPI.register(payload);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem('token');
    setToken(null);
    setUser(null);
  }, []);

  const isAdmin = user?.role === 'admin';
  const isAuthenticated = !!token;

  return (
    <AuthContext.Provider value={{ user, token, isAuthenticated, isAdmin, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth() {
  return useContext(AuthContext);
}
