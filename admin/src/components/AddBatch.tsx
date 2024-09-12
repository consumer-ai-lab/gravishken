import React, { useState } from 'react';
import { Batch } from '@common/types';

const AddBatch: React.FC = () => {
  const [batchName, setBatchName] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const batch: Batch = {
      batchName,
    };

    try {
      const response = await fetch('http://localhost:8081/admin/add_batch', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(batch),
      });

      if (response.ok) {
        alert('Batch added successfully!');
        setBatchName('');
      } else {
        throw new Error('Failed to add batch');
      }
    } catch (error) {
      console.error('Error adding batch:', error);
      alert('Error adding batch');
    }
  };

  return (
    <div className="max-w-lg mx-auto p-6 bg-gray-100 rounded-lg shadow-md">
      <h2 className="text-2xl font-semibold text-center mb-6">Add New Batch</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="batchName" className="block text-sm font-medium text-gray-700">
            Batch Name:
          </label>
          <input
            type="text"
            id="batchName"
            value={batchName}
            onChange={(e) => setBatchName(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            required
          />
        </div>

        <div>
          <button
            type="submit"
            className="w-full px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-md shadow hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            Add Batch
          </button>
        </div>
      </form>
    </div>
  );
};

export default AddBatch;
