import React, { useState } from 'react';
import { TestType, Test } from '@common/types';

const AddTest: React.FC = () => {
  const [tests, setTests] = useState<Test[]>([]);
  const [testDuration, setTestDuration] = useState<number>(60);
  const [password, setPassword] = useState('');
  const [batchId, setBatchId] = useState('');
  const [currentTest, setCurrentTest] = useState<Test>();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const batchTest = {
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


  return (
    <div>
      Working.
    </div>
  );
};

export default AddTest;
