
"use client"
import React, { useState } from 'react';

interface BatchData {
  id: string;
  batchName: string;
  [key: string]: any;
}

const GetBatchwiseData: React.FC = () => {
  const [param, setParam] = useState<string>('roll');
  const [batchNumber, setBatchNumber] = useState<string>('1');
  const [rangeStart, setRangeStart] = useState<number>(12767);
  const [rangeEnd, setRangeEnd] = useState<number>(12928);
  const [batchData, setBatchData] = useState<BatchData[] | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    setLoading(true);
    setErrorMessage(null);

    const requestBody = {
      param,
      batchNumber,
      ranges: [rangeStart, rangeEnd],
    };

    try {
      const response = await fetch('http://localhost:8081/admin/get_batchwise_data', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestBody),
      });

      if (response.ok) {
        const data = await response.json();
        // Update this line to extract the `data` field from the response
        // Access the "data" field in the response
        if (Array.isArray(data.data)) {
          setBatchData(data.data); // Set the batches to the array in "data"
        } else {
            throw new Error('Unexpected data format');
        }
      } else {
        throw new Error('Failed to fetch batch data');
      }
    } catch (err) {
      setErrorMessage('Error fetching batch data');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Get Batchwise Data</h1>

      {/* Form to input the parameters */}
      <form onSubmit={handleSubmit}>
        {/* Param Dropdown */}
        <div className="mb-4">
          <label htmlFor="param" className="block text-gray-700 font-bold mb-2">
            Parameter:
          </label>
          <select
            id="param"
            value={param}
            onChange={(e) => setParam(e.target.value)}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          >
            <option value="roll">Roll</option>
            <option value="batch">Batch</option>
            <option value="frontend">Frontend</option>
          </select>
        </div>

        {/* Batch Number input */}
        <div className="mb-4">
          <label htmlFor="batchNumber" className="block text-gray-700 font-bold mb-2">
            Batch Number:
          </label>
          <input
            id="batchNumber"
            type="text"
            value={batchNumber}
            onChange={(e) => setBatchNumber(e.target.value)}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          />
        </div>

        {/* Ranges inputs */}
        <div className="mb-4">
          <label htmlFor="rangeStart" className="block text-gray-700 font-bold mb-2">
            Range Start:
          </label>
          <input
            id="rangeStart"
            type="number"
            value={rangeStart}
            onChange={(e) => setRangeStart(Number(e.target.value))}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          />
        </div>
        <div className="mb-4">
          <label htmlFor="rangeEnd" className="block text-gray-700 font-bold mb-2">
            Range End:
          </label>
          <input
            id="rangeEnd"
            type="number"
            value={rangeEnd}
            onChange={(e) => setRangeEnd(Number(e.target.value))}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          />
        </div>

        {/* Submit button */}
        <button
          type="submit"
          className={`bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
          disabled={loading}
        >
          {loading ? 'Fetching...' : 'Get Batch Data'}
        </button>
      </form>

      {/* Display Batch Data */}
      {errorMessage && <p className="text-red-500 mt-4">{errorMessage}</p>}
      {batchData && (
        <div className="mt-6">
          <h2 className="text-xl font-bold">Batch Data:</h2>
          <ul className="mt-4">
            {batchData.map((batch: BatchData, index: number) => (
              <li key={index} className="mb-2 p-2 border border-gray-300 rounded">
                <p><strong>ID:</strong> {batch.id}</p>
                <p><strong>Batch Name:</strong> {batch.batchName}</p>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default GetBatchwiseData;
