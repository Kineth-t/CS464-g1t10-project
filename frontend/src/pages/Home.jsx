import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { phonesAPI } from '../api/client';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Search, Package } from 'lucide-react';

export default function Home() {
  const [phones, setPhones] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');

  useEffect(() => {
    phonesAPI.list()
      .then((data) => setPhones(Array.isArray(data) ? data : []))
      .catch((e) => {
        console.error('Failed to load phones:', e);
        setError(e.message);
      })
      .finally(() => setLoading(false));
  }, []);

  const filtered = phones.filter((p) => {
    const q = search.toLowerCase();
    return p.brand.toLowerCase().includes(q) || p.model.toLowerCase().includes(q);
  });

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <div className="flex flex-col gap-2 mb-8">
        <h1 className="text-3xl font-bold tracking-tight">Browse Phones</h1>
        <p className="text-muted-foreground">Find your next device from our collection</p>
      </div>

      <div className="relative mb-6 max-w-sm">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          className="pl-9"
          placeholder="Search by brand or model..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>

      {loading && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {[...Array(8)].map((_, i) => (
            <Card key={i} className="animate-pulse">
              <div className="h-48 bg-muted rounded-t-lg" />
              <CardHeader><div className="h-5 bg-muted rounded w-3/4" /></CardHeader>
              <CardContent><div className="h-4 bg-muted rounded w-1/2" /></CardContent>
            </Card>
          ))}
        </div>
      )}

      {error && (
        <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-4 text-destructive text-sm">
          {error}
        </div>
      )}

      {!loading && !error && filtered.length === 0 && (
        <div className="flex flex-col items-center gap-3 py-16 text-muted-foreground">
          <Package className="h-12 w-12 opacity-40" />
          <p className="text-lg font-medium">No phones found</p>
          <p className="text-sm">{search ? 'Try a different search term' : 'Check back later for new arrivals'}</p>
        </div>
      )}

      {!loading && !error && filtered.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {filtered.map((phone) => (
            <Link to={`/phones/${phone.id}`} key={phone.id} className="group">
              <Card className="h-full flex flex-col transition-shadow hover:shadow-md">

                <div className="flex items-center justify-center bg-muted/30 rounded-t-lg h-48 px-4 pt-4">
                  {phone.image_url ? (
                    <img
                      src={phone.image_url}
                      alt={`${phone.brand} ${phone.model}`}
                      className="h-full w-full object-contain drop-shadow-sm"
                    />
                  ) : (
                    <Package className="h-16 w-16 text-muted-foreground opacity-30" />
                  )}
                </div>

                <CardHeader className="pb-2">
                  <CardTitle className="text-base">{phone.brand} {phone.model}</CardTitle>
                </CardHeader>

                <CardContent className="flex-1 pb-3">
                  <p className="text-2xl font-bold text-primary">S${Number(phone.price).toFixed(2)}</p>
                  <div className="mt-2">
                    {phone.stock > 0
                      ? <Badge variant="secondary">{phone.stock} in stock</Badge>
                      : <Badge variant="destructive">Out of stock</Badge>
                    }
                  </div>
                  {phone.description && (
                    <p className="mt-2 text-sm text-muted-foreground line-clamp-2">{phone.description}</p>
                  )}
                </CardContent>

                <CardFooter>
                  <Button variant="outline" size="sm" className="w-full group-hover:bg-primary group-hover:text-primary-foreground transition-colors">
                    View Details
                  </Button>
                </CardFooter>

              </Card>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}