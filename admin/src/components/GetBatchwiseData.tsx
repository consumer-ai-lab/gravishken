import React, { useState } from 'react';

const GetBatchwiseData: React.FC = () => {
  const [batchName, setBatchName] = useState('');
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setData(null);

    try {
      const response = await fetch(`http://localhost:8081/admin/get_batchwise_data?batch=${batchName}`);
      if (!response.ok) {
        throw new Error('Failed to fetch batchwise data');
      }
      const result = await response.json();
      setData(result);
    } catch (err) {
      setError('Error fetching batchwise data');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Get Batchwise Data</h1>
      <form onSubmit={handleSubmit} className="mb-4">
        <div className="mb-4">
          <label htmlFor="batchName" className="block text-gray-700 font-bold mb-2">
            Batch Name:
          </label>
          <input
            id="batchName"
            type="text"
            value={batchName}
            onChange={(e) => setBatchName(e.target.value)}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            required
          />
        </div>
        <button
          type="submit"
          className={`bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
          disabled={loading}
        >
          {loading ? 'Fetching...' : 'Get Data'}
        </button>
      </form>
      {error && <p className="text-red-500 mb-4">{error}</p>}
      {data && (
        <div>
          <h2 className="text-xl font-bold mb-2">Batch Data:</h2>
          <pre className="bg-gray-100 p-4 rounded overflow-auto">
            {JSON.stringify(data, null, 2)}
          </pre>
        </div>
      )}
    </div>
  );
};

export default GetBatchwiseData;