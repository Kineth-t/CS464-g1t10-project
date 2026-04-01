import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Smartphone, AlertCircle } from 'lucide-react';

const INITIAL = {
  username: '',
  password: '',
  phone_number: '',
  street: '',
  city: '',
  state: '',
  country: '',
  zip_code: '',
};

function Field({ label, name, type = 'text', required = false, value, onChange }) {
  return (
    <div className="space-y-2">
      <Label htmlFor={name}>
        {label}
        {required && <span className="text-destructive ml-1">*</span>}
      </Label>
      <Input
        id={name}
        name={name}
        type={type}
        value={value}
        onChange={onChange}
        required={required}
      />
    </div>
  );
}

export default function Register() {
  const { register } = useAuth();
  const navigate = useNavigate();

  const [form, setForm] = useState(INITIAL);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  function handleChange(e) {
    setForm((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await register(form);
      navigate('/login');
    } catch (err) {
      setError(err.message || 'Something went wrong');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-[calc(100vh-56px)] flex items-center justify-center px-4 py-12 bg-muted/30">
      <div className="w-full max-w-lg">
        <div className="flex justify-center mb-6">
          <div className="flex items-center gap-2 font-bold text-xl">
            <Smartphone className="h-6 w-6" /> Ringr Mobile
          </div>
        </div>

        <Card>
          <CardHeader className="space-y-1">
            <CardTitle className="text-2xl">Create an account</CardTitle>
            <CardDescription>
              Fill in your details to get started
            </CardDescription>
          </CardHeader>

          <form onSubmit={handleSubmit}>
            <CardContent className="space-y-4">
              {error && (
                <div className="flex items-center gap-2 rounded-md bg-destructive/10 text-destructive px-3 py-2 text-sm">
                  <AlertCircle className="h-4 w-4 shrink-0" />
                  {error}
                </div>
              )}

              <div className="grid grid-cols-2 gap-4">
                <Field
                  label="Username"
                  name="username"
                  value={form.username}
                  onChange={handleChange}
                  required
                />
                <Field
                  label="Password"
                  name="password"
                  type="password"
                  value={form.password}
                  onChange={handleChange}
                  required
                />
              </div>

              <Field
                label="Phone Number"
                name="phone_number"
                value={form.phone_number}
                onChange={handleChange}
              />

              <Separator />

              <p className="text-sm font-medium text-muted-foreground">
                Delivery Address
              </p>

              <Field
                label="Street"
                name="street"
                value={form.street}
                onChange={handleChange}
              />

              <div className="grid grid-cols-2 gap-4">
                <Field
                  label="City"
                  name="city"
                  value={form.city}
                  onChange={handleChange}
                />
                <Field
                  label="State"
                  name="state"
                  value={form.state}
                  onChange={handleChange}
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <Field
                  label="Country"
                  name="country"
                  value={form.country}
                  onChange={handleChange}
                />
                <Field
                  label="Zip Code"
                  name="zip_code"
                  value={form.zip_code}
                  onChange={handleChange}
                />
              </div>
            </CardContent>

            <CardFooter className="flex flex-col gap-3">
              <Button
                className="w-full"
                type="submit"
                disabled={loading}
              >
                {loading ? 'Creating account...' : 'Create account'}
              </Button>

              <p className="text-sm text-muted-foreground text-center">
                Already have an account?{' '}
                <Link
                  to="/login"
                  className="text-primary underline-offset-4 hover:underline"
                >
                  Sign in
                </Link>
              </p>
            </CardFooter>
          </form>
        </Card>
      </div>
    </div>
  );
}