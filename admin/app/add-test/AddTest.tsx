"use client";

import { useState } from 'react';

const AddTest = () => {
  const [fileType, setFileType] = useState('pdf');
  const [timeSlot, setTimeSlot] = useState('');
  const [password, setPassword] = useState('');
  const [batch, setBatch] = useState('');
  const [driveId, setDriveId] = useState('');

  const handleSubmit = (e:any) => {
    e.preventDefault();

    // Validate drive link format
    const driveRegex = /^(https:\/\/)?(www\.)?drive\.google\.com\/.*$/;
    if (!driveRegex.test(driveId)) {
      alert('Please enter a valid Google Drive link');
      return;
    }

    const testDetails = {
      fileType,
      timeSlot,
      password,
      batch,
      driveId,
    };

    console.log('Test Details:', testDetails);
    alert('Test added successfully!');
  };

  return (
    <div className="max-w-lg mx-auto p-6 bg-gray-100 rounded-lg shadow-md">
      <h2 className="text-2xl font-semibold text-center mb-6">Add New Test</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        
        {/* File Type Selection */}
        <div>
          <label htmlFor="fileType" className="block text-sm font-medium text-gray-700">
            File Type:
          </label>
          <select
            id="fileType"
            value={fileType}
            onChange={(e) => setFileType(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          >
            <option value="pdf">PDF</option>
            <option value="doc">DOC</option>
            <option value="txt">TXT</option>
          </select>
        </div>

        {/* Time Slot Input */}
        <div>
          <label htmlFor="timeSlot" className="block text-sm font-medium text-gray-700">
            Time Slot:
          </label>
          <input
            type="datetime-local"
            id="timeSlot"
            value={timeSlot}
            onChange={(e) => setTimeSlot(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            required
          />
        </div>

        {/* Password Input */}
        <div>
          <label htmlFor="password" className="block text-sm font-medium text-gray-700">
            Password:
          </label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            required
          />
        </div>

        {/* Batch Input */}
        <div>
          <label htmlFor="batch" className="block text-sm font-medium text-gray-700">
            Batch:
          </label>
          <input
            type="number"
            id="batch"
            value={batch}
            onChange={(e) => setBatch(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            required
            min="1"
          />
        </div>

        {/* Drive Link Input */}
        <div>
          <label htmlFor="driveId" className="block text-sm font-medium text-gray-700">
            Drive Link:
          </label>
          <input
            type="url"
            id="driveId"
            value={driveId}
            onChange={(e) => setDriveId(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            required
            placeholder="https://drive.google.com/..."
          />
        </div>

        {/* Submit Button */}
        <div>
          <button
            type="submit"
            className="w-full px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-md shadow hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            Submit
          </button>
        </div>
      </form>
    </div>
  );
};

export default AddTest;
