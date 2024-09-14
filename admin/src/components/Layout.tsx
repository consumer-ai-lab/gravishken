import React from 'react';
import Navbar from './Navbar';
import Cookies from "js-cookie"
import { Navigate } from 'react-router-dom';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {

  const hasAdminCookie = Cookies.get('admin_data') !== undefined;

  if (!hasAdminCookie) {
    return <Navigate to="/login" replace />;
  }

  return (
    <div>
      <Navbar />
      <main className="container mx-auto mt-8 px-4">
        {children}
      </main>
    </div>
  );
};

export default Layout;
