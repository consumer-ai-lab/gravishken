"use client";

import React, { useState } from 'react';

const UpdateUserData: React.FC = () => {
  const [username, setUsername] = useState<string>('');
  const [property, setProperty] = useState<string>('start_time');
  const [value, setValue] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    setLoading(true);
    setErrorMessage(null);
    setSuccessMessage(null);

    const requestBody = {
      username,
      property,
      value: [value], // Wrapping value in an array, as shown in the request example
    };

    try {
      const response = await fetch('http://localhost:8081/admin/update_user_data', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestBody),
      });

      if (response.ok) {
        setSuccessMessage('User data updated successfully!');
      } else {
        throw new Error('Failed to update user data');
      }
    } catch (err) {
      setErrorMessage('Error updating user data');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Update User Data</h1>

      {/* Form for updating user data */}
      <form onSubmit={handleSubmit}>
        {/* Username input */}
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

        {/* Property dropdown */}
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
            <option value="start_time">Start Time</option>
            <option value="reading_submission_received">Reading Submission Received</option>
            <option value="submission_received">Submission Received</option>
            <option value="elapsed_time">Elapsed Time</option>
            <option value="reading_elapsed_time">Reading Elapsed Time</option>
            <option value="submission_folder_id">Submission Folder ID</option>
            <option value="wpm">WPM</option>
            <option value="user_test_time">User Test Time</option>
            <option value="batch_test_time">Batch Test Time</option>
          </select>
        </div>

        {/* Value input */}
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

        {/* Submit button */}
        <button
          type="submit"
          className={`bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
          disabled={loading}
        >
          {loading ? 'Updating...' : 'Update User Data'}
        </button>
      </form>

      {/* Success and Error Messages */}
      {successMessage && <p className="text-green-500 mt-4">{successMessage}</p>}
      {errorMessage && <p className="text-red-500 mt-4">{errorMessage}</p>}
    </div>
  );
};

export default UpdateUserData;
