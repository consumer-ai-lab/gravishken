import React, { useState } from 'react';
import { Link, Outlet } from 'react-router-dom';

const Dashboard: React.FC = () => {
  const [activeFeature, setActiveFeature] = useState('');

  const features = [
    { name: "Add Test", to: "/add-test" },
    { name: "Add Batch", to: "/add-batch" },
    { name: "Add All Users", to: "/add-users" },
    { name: "Update User Data", to: "/update-user" },
    { name: "Get Batchwise Data", to: "/get-batchwise-data" },
    { name: "Increase Test Time", to: "/increase-test-time" },
    { name: "Set User Data", to: "/set-user-data" },
  ];

  return (
    <div className="flex h-screen">
      {/* Left part - Navigation */}
      <div className="w-1/4 bg-gray-100 p-4 overflow-y-auto">
        <h2 className="text-2xl font-bold mb-4">Admin Dashboard</h2>
        <nav>
          {features.map((feature, index) => (
            <Link
              key={index}
              to={feature.to}
              className={`block mb-2 p-2 rounded ${activeFeature === feature.to ? 'bg-blue-500 text-white' : 'hover:bg-gray-200'}`}
              onClick={() => setActiveFeature(feature.to)}
            >
              {feature.name}
            </Link>
          ))}
        </nav>
      </div>

      {/* Right part - Content */}
      <div className="w-3/4 p-8 overflow-y-auto">
        <Outlet />
      </div>
    </div>
  );
};

export default Dashboard;
