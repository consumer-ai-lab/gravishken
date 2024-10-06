"use client";

import React, { useState, useEffect } from 'react';

interface Batch {
    merged_file_id: string,
    resultDownloaded: boolean,
    submission_folder_id: string,
    submission_received: boolean,
    username: string
}

const AddBatch: React.FC = () => {
    const [batches, setBatches] = useState<Batch[]>([]);
    const [batchName, setBatchName] = useState<string>('');
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [successMessage, setSuccessMessage] = useState<string | null>(null);
  
    // Fetch available batches
    useEffect(() => {
    const fetchBatches = async () => {
        try {
        const response = await fetch('http://localhost:6201/batch/get_batches');
        if (!response.ok) {
            throw new Error('Failed to fetch batches');
        }
        const data = await response.json();

        // Access the "data" field in the response
        if (Array.isArray(data.data)) {
            setBatches(data.data); // Set the batches to the array in "data"
        } else {
            throw new Error('Unexpected data format');
        }
        } catch (err) {
        setError('Failed to load batches');
        } finally {
        setLoading(false);
        }
    };

    fetchBatches();
    }, []);

    console.log(batches);

  const handleBatchNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setBatchName(e.target.value);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const response = await fetch('http://localhost:6201/admin/add_batch', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ batchName }),
      });

      if (response.ok) {
        setSuccessMessage('Batch added successfully!');
        setBatchName(''); // Clear the input field after successful submission
        // Re-fetch batches to include the new one
        const fetchBatches = async () => {
          const updatedResponse = await fetch('http://localhost:6201/batch/get_batches');
          const updatedData = await updatedResponse.json();
          setBatches(updatedData.data); // Update the batches list
        };
        fetchBatches();
      } else {
        throw new Error('Failed to add batch');
      }
    } catch (err) {
      setError('Error adding batch');
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Add Batch</h1>

      {/* Form to add a new batch */}
      <form onSubmit={handleSubmit} className="mb-6">
        <div className="mb-4">
          <label htmlFor="batchName" className="block text-gray-700 font-bold mb-2">
            Batch Name:
          </label>
          <input
            id="batchName"
            type="text"
            value={batchName}
            onChange={handleBatchNameChange}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            placeholder="Enter batch name"
            required
          />
        </div>
        <button
          type="submit"
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
        >
          Add Batch
        </button>
      </form>

      {/* Loading state */}
      {loading ? (
        <div className="flex justify-center items-center">
          <p className="text-lg">Loading batches...</p>
        </div>
      ) : error ? (
        <div className="text-red-500 text-center">{error}</div>
      ) : (
        <div>
          <h2 className="text-xl font-bold mb-4">Available Batches</h2>
          <ul className="list-disc list-inside">
            {batches.map((batch) => (
              <ul key={batch.merged_file_id} className="mb-2">
                <li >
                    {batch.merged_file_id}
                </li>
                <li>
                    {batch.username}
                </li>
                <li>
                    {batch.submission_folder_id}
                </li>
                <li>
                    {batch.resultDownloaded}
                </li>
                <li>
                    {batch.submission_received}
                </li>
              </ul>
              
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default AddBatch;
