import React, { useState } from 'react';

const IncreaseTestTime: React.FC = () => {
  const [testPassword, setTestPassword] = useState('');
  const [additionalTime, setAdditionalTime] = useState('');
  const [message, setMessage] = useState('');

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    // Implement the logic to increase test time here
    console.log({ testPassword, additionalTime });
    setMessage('Test time increased successfully!');
  };

  return (
    <div className="max-w-md mx-auto mt-10">
      <h2 className="text-2xl font-bold mb-4">Increase Test Time</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="testPassword" className="block mb-1">
            Test Password:
          </label>
          <input
            type="password"
            id="testPassword"
            value={testPassword}
            onChange={(e) => setTestPassword(e.target.value)}
            required
            className="w-full px-3 py-2 border rounded-md"
          />
        </div>
        <div>
          <label htmlFor="additionalTime" className="block mb-1">
            Additional Time (in minutes):
          </label>
          <input
            type="number"
            id="additionalTime"
            value={additionalTime}
            onChange={(e) => setAdditionalTime(e.target.value)}
            required
            className="w-full px-3 py-2 border rounded-md"
          />
        </div>
        <button
          type="submit"
          className="w-full bg-green-500 text-white py-2 rounded-md hover:bg-green-600"
        >
          Increase Time
        </button>
      </form>
      {message && (
        <div className="mt-4 p-2 bg-green-100 text-green-700 rounded">
          {message}
        </div>
      )}
    </div>
  );
};

export default IncreaseTestTime;
