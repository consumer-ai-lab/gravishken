"use client";

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
      <div className="container mx-auto flex gap-y-3 justify-between items-center">
        <div className="text-xl font-semibold text-white">
          <a href="/adminDashboard">Admin Dashboard</a>
        </div>
        <div className="flex-grow flex space-x-6">
          <a href="/adminViewSubmission" className="text-white hover:text-gray-300">
            Download All Files
          </a>
          <a href="/adminLoginStatus" className="text-white hover:text-gray-300">
            Check Login Status
          </a>
          <a href="/increaseTime" className="text-white hover:text-gray-300">
            Increase Time
          </a>
          <a href="/resetValues" className="text-white hover:text-gray-300">
            Reset Values
          </a>
          <a href="/increaseBatchTime" className="text-white hover:text-gray-300">
            Increase Batch Time
          </a>
          <a href="/adminWPMResults" className="text-white hover:text-gray-300">
            WPM Results
          </a>
          <a href="/downloadRoll" className="text-white hover:text-gray-300">
            Download Rolls
          </a>
        </div>
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2 text-white">
            <span>Welcome</span>
            <span>{"admin"}</span>
          </div>
          <button
            onClick={handleLogout}
            className="flex items-center bg-blue-600 text-white px-3 py-2 rounded hover:bg-blue-700 transition"
          >
            Login
          </button>
        </div>
      </div>
    </nav>
  );
};

export default AdminNavbar;
