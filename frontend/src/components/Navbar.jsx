import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { buttonVariants } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { ShoppingCart, Shield, Smartphone, ChevronDown, LogOut } from 'lucide-react';
import { cn } from '@/lib/utils';

export default function Navbar() {
  const { isAuthenticated, isAdmin, user, logout } = useAuth();
  const navigate = useNavigate();

  function handleLogout() {
    logout();
    navigate('/login');
  }

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="max-w-6xl mx-auto px-4 h-14 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2 font-bold text-lg">
          <Smartphone className="h-5 w-5" />
          Ringr Mobile
        </Link>

        <nav className="flex items-center gap-1">
          <Link to="/" className={cn(buttonVariants({ variant: 'ghost', size: 'sm' }))}>
            Phones
          </Link>

          {isAuthenticated && (
            <Link to="/cart" className={cn(buttonVariants({ variant: 'ghost', size: 'sm' }), 'gap-1')}>
              <ShoppingCart className="h-4 w-4" /> Cart
            </Link>
          )}

          {isAdmin && (
            <Link to="/admin" className={cn(buttonVariants({ variant: 'ghost', size: 'sm' }), 'gap-1')}>
              <Shield className="h-4 w-4" /> Admin
            </Link>
          )}

          <Separator orientation="vertical" className="h-6 mx-2" />

          {isAuthenticated ? (
            <DropdownMenu>
              <DropdownMenuTrigger className={cn(buttonVariants({ variant: 'ghost', size: 'sm' }), 'gap-2')}>
                <Avatar className="h-6 w-6">
                  <AvatarFallback className="text-xs">
                    {user?.username?.[0]?.toUpperCase()}
                  </AvatarFallback>
                </Avatar>
                {user?.username}
                <ChevronDown className="h-3 w-3 opacity-50" />
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={handleLogout} className="text-destructive gap-2 cursor-pointer">
                  <LogOut className="h-4 w-4" /> Logout
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : (
            <div className="flex items-center gap-2">
              <Link to="/login" className={cn(buttonVariants({ variant: 'ghost', size: 'sm' }))}>
                Login
              </Link>
              <Link to="/register" className={cn(buttonVariants({ size: 'sm' }))}>
                Register
              </Link>
            </div>
          )}
        </nav>
      </div>
    </header>
  );
}
