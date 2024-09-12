import React from 'react';
import { Link } from 'react-router-dom';

const Navbar: React.FC = () => {
  return (
    <nav className="bg-gray-800 text-white p-4">
      <div className="max-w-7xl mx-auto flex justify-between items-center">
        <Link to="/" className="text-xl font-bold">Admin Dashboard</Link>
        <Link to="/" className="hover:text-gray-300">Home</Link>
      </div>
    </nav>
  );
};

export default Navbar;
