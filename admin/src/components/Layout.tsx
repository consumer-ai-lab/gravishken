import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Navigate, useNavigate } from 'react-router-dom';
import { useToast } from "@/hooks/use-toast"

interface LayoutProps {
  children: React.ReactNode;
}

export default function Layout({ children }:LayoutProps) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
  const navigate = useNavigate();
  const { toast } = useToast();

  useEffect(() => {
    const checkAuthStatus = async () => {
      try {
        const response = await axios.get(`${import.meta.env.VITE_BACKEND_BASE_URL}/admin/auth-status`, {
          withCredentials: true,
        });

        
        if (response.data.isAuthenticated) {
          setIsAuthenticated(true);
        } else {
          setIsAuthenticated(false);
          navigate('/login');
        }
      } catch (error) {
        console.error('Authentication check failed:', error);
        setIsAuthenticated(false);
        toast({
          variant: "destructive",
          title: "Authentication Failed",
          description: "Please log in again.",
        });
        navigate('/login');
      }
    };

    checkAuthStatus();
  }, [navigate, toast]);



 if (isAuthenticated === null) {
    return <div>Loading...</div>;
  }

  if (isAuthenticated === false) {
    return <Navigate to="/login" replace />;
  }

  return (
    <div>
      <main className="container mx-auto mt-8 px-4">
        {children}
      </main>
    </div>
  );
};
