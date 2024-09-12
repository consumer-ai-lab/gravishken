import React, { useState } from 'react';
import { User } from '@common/types';

const AddUser: React.FC = () => {
  const [user, setUser] = useState<User>({
    name: '',
    username: '',
    password: '',
    testPassword: '',
    batch: '',
    id: undefined,
    tests: undefined,
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setUser(prevUser => ({ ...prevUser, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await fetch('http://localhost:8081/admin/add_user', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(user),
      });

      if (response.ok) {
        alert('User added successfully!');
        setUser({ name: '', username: '', password: '', testPassword: '', batch: '' });
      } else {
        throw new Error('Failed to add user');
      }
    } catch (error) {
      console.error('Error adding user:', error);
      alert('Error adding user');
    }
  };

  return (
    <div className="max-w-lg mx-auto p-6 bg-gray-100 rounded-lg shadow-md">
      <h2 className="text-2xl font-semibold text-center mb-6">Add New User</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        {Object.keys(user).map((key) => (
          <div key={key}>
            <label htmlFor={key} className="block text-sm font-medium text-gray-700">
              {key.charAt(0).toUpperCase() + key.slice(1)}:
            </label>
            <input
              type={key.includes('password') ? 'password' : 'text'}
              id={key}
              name={key}
              // @ts-ignore
              value={user[key as keyof User]}
              onChange={handleChange}
              className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              required
            />
          </div>
        ))}
        <button
          type="submit"
          className="w-full px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-md shadow hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          Add User
        </button>
      </form>
    </div>
  );
};

export default AddUser;