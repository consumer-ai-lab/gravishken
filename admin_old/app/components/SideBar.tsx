import React from 'react';
import Link from 'next/link';

const SideBar: React.FC = () => {
  const features = [
    { name: "Add Test", path: "/add-test" },
    { name: "Add Batch", path: "/add-batch" },
    { name: "Add All Users", path: "/add-users" },
    { name: "Update User Data", path: "/update-user" },
    { name: "Get Batchwise Data", path: "/get-batchwise-data" },
    { name: "Increase Test Time", path: "/increase-test-time" },
    { name: "Set User Data", path: "/set-user-data" }
  ];

  return (
    <div className='flex flex-col space-y-4'>
      {features.map((feature, index) => (
        <Link key={index} href={feature.path} passHref>
          <button className='m-2 p-4 bg-gray-200 rounded-md cursor-pointer hover:bg-gray-300 transition'>
            {feature.name}
          </button>
        </Link>
      ))}
    </div>
  );
};

export default SideBar;
