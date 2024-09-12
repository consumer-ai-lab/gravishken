"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

const AdminNavbar: React.FC<any> = () => {
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) return null; 

  const handleLogout = () => {
    console.log("Logging out");
  };

  return (
    <nav className="bg-gray-900 w-full shadow-lg p-4">
      <Link href={"/"} className="text-white font-medium">
        Admin Panel
      </Link>
    </nav>
  );
};

export default AdminNavbar;
