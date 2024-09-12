import React, { useState } from 'react';
import { TestType, Test, BatchTests, FileType } from '@common/types';

const AddTest: React.FC = () => {
  const [tests, setTests] = useState<Test[]>([]);
  const [testDuration, setTestDuration] = useState<number>(60);
  const [password, setPassword] = useState('');
  const [batchId, setBatchId] = useState('');
  const [currentTest, setCurrentTest] = useState<Test>({
    testId: [],
    testType: TestType.FileTest,
    fileType: FileType.DOCX,
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const batchTest: BatchTests = {
      batchId: [], // This will be filled on the server
      tests: tests,
      testDuration,
      password,
      startTime: new Date(),
      endTime: new Date(Date.now() + testDuration * 60000),
    };

    try {
      const response = await fetch('http://localhost:8081/admin/add_test', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(batchTest),
      });

      if (response.ok) {
        alert('Tests added successfully!');
        setTests([]);
      } else {
        throw new Error('Failed to add tests');
      }
    } catch (error) {
      console.error('Error adding tests:', error);
      alert('Error adding tests');
    }
  };

  const addTestToList = () => {
    setTests([...tests, currentTest]);
    setCurrentTest({
      testId: [],
      testType: TestType.FileTest,
      fileType: FileType.DOCX,
    });
  };

  return (
    <div className="max-w-lg mx-auto p-6 bg-gray-100 rounded-lg shadow-md">
      <h2 className="text-2xl font-semibold text-center mb-6">Add New Batch Test</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="testDuration" className="block text-sm font-medium text-gray-700">
            Total Test Duration (minutes):
          </label>
          <input
            type="number"
            id="testDuration"
            value={testDuration}
            onChange={(e) => setTestDuration(parseInt(e.target.value))}
            className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            required
            min="1"
          />
        </div>

        <div>
          <label htmlFor="password" className="block text-sm font-medium text-gray-700">
            Test Code:
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

        <div>
          <label htmlFor="batchId" className="block text-sm font-medium text-gray-700">
            Batch ID:
          </label>
          <input
            type="text"
            id="batchId"
            value={batchId}
            onChange={(e) => setBatchId(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            required
          />
        </div>

        <div className="border-t pt-4">
          <h3 className="text-lg font-semibold mb-2">Add Individual Tests</h3>
          <div>
            <label htmlFor="testType" className="block text-sm font-medium text-gray-700">
              Test Type:
            </label>
            <select
              id="testType"
              value={currentTest.testType}
              onChange={(e) => setCurrentTest({...currentTest, testType: e.target.value as TestType})}
              className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            >
              {Object.values(TestType).map((type) => (
                <option key={type} value={type}>
                  {type}
                </option>
              ))}
            </select>
          </div>

          {currentTest.testType === TestType.FileTest && (
            <div>
              <label htmlFor="fileType" className="block text-sm font-medium text-gray-700">
                File Type:
              </label>
              <select
                id="fileType"
                value={currentTest.fileType}
                onChange={(e) => setCurrentTest({...currentTest, fileType: e.target.value as FileType})}
                className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              >
                {Object.values(FileType).map((type) => (
                  <option key={type} value={type}>
                    {type}
                  </option>
                ))}
              </select>
            </div>
          )}

          {currentTest.testType === TestType.TypingTest && (
            <div>
              <label htmlFor="typingTestText" className="block text-sm font-medium text-gray-700">
                Typing Test Text:
              </label>
              <textarea
                id="typingTestText"
                value={currentTest.typingTestText || ''}
                onChange={(e) => setCurrentTest({...currentTest, typingTestText: e.target.value})}
                className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                rows={6}
                required={currentTest.testType === TestType.TypingTest}
              />
            </div>
          )}

          <button
            type="button"
            onClick={addTestToList}
            className="mt-2 px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md shadow hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
          >
            Add Test
          </button>
        </div>

        {tests.length > 0 && (
          <div>
            <h3 className="text-lg font-semibold mb-2">Current Tests:</h3>
            <ul className="space-y-2">
              {tests.map((test, index) => (
                <li key={index} className="bg-white p-3 rounded-md shadow">
                  <p><strong>Test Type:</strong> {test.testType}</p>
                  {test.testType === TestType.FileTest && (
                    <p><strong>File Type:</strong> {test.fileType}</p>
                  )}
                  {test.testType === TestType.TypingTest && (
                    <p><strong>Typing Test Text:</strong> {test.typingTestText?.substring(0, 50)}...</p>
                  )}
                </li>
              ))}
            </ul>
          </div>
        )}

        <div>
          <button
            type="submit"
            className="w-full px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-md shadow hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            Submit Test Batch
          </button>
        </div>
      </form>
    </div>
  );
};

export default AddTest;
