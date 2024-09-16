import React, { useState } from 'react';
import { UserUpdateRequest } from '@common/types';

const UpdateUserData: React.FC = () => {
  const [username, setUsername] = useState('');
  const [property, setProperty] = useState('');
  const [value, setValue] = useState('');
  const [loading, setLoading] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setSuccessMessage('');
    setErrorMessage('');

    const updateRequest: UserUpdateRequest = {
      username,
      property,
      value: [value], // Assuming the value is always an array of strings
    };

    try {
      const response = await fetch(`${import.meta.env.SERVER_URL}/admin/update_user_data`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updateRequest),
      });

      if (response.ok) {
        setSuccessMessage('User data updated successfully!');
      } else {
        throw new Error('Failed to update user data');
      }
    } catch (error) {
      console.error('Error updating user data:', error);
      setErrorMessage('Error updating user data');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Update User Data</h1>
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label htmlFor="username" className="block text-gray-700 font-bold mb-2">
            Username:
          </label>
          <input
            id="username"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            required
          />
        </div>
        <div className="mb-4">
          <label htmlFor="property" className="block text-gray-700 font-bold mb-2">
            Property to Update:
          </label>
          <select
            id="property"
            value={property}
            onChange={(e) => setProperty(e.target.value)}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            required
          >
            <option value="">Select a property</option>
            <option value="name">Name</option>
            <option value="password">Password</option>
            <option value="testPassword">Test Password</option>
            <option value="batch">Batch</option>
          </select>
        </div>
        <div className="mb-4">
          <label htmlFor="value" className="block text-gray-700 font-bold mb-2">
            Value:
          </label>
          <input
            id="value"
            type="text"
            value={value}
            onChange={(e) => setValue(e.target.value)}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            required
          />
        </div>
        <button
          type="submit"
          className={`bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
          disabled={loading}
        >
          {loading ? 'Updating...' : 'Update User Data'}
        </button>
      </form>
      {successMessage && <p className="text-green-500 mt-4">{successMessage}</p>}
      {errorMessage && <p className="text-red-500 mt-4">{errorMessage}</p>}
    </div>
  );
};

export default UpdateUserData;
