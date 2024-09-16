import React, { useState } from 'react';

const UpdateTypingTestText: React.FC = () => {
  const [testPassword, setTestPassword] = useState('');
  const [typingTestText, setTypingTestText] = useState('');
  const [message, setMessage] = useState('');

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const response = await fetch(`${import.meta.env.SERVER_URL}/admin/update_typing_test_text`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ testPassword, typingTestText }),
      });

      if (response.ok) {
        setMessage('Typing test updated successfully!');
      } else {
        setMessage('Failed to update typing test. Please try again.');
      }
    } catch (error) {
      setMessage('An error occurred. Please try again later.');
    }
  };

  return (
    <div className="max-w-md mx-auto mt-10">
      <h2 className="text-2xl font-bold mb-4">Update Typing Test</h2>
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
          <label htmlFor="typingTestText" className="block mb-1">
            Typing Test Text:
          </label>
          <textarea
            id="typingTestText"
            value={typingTestText}
            onChange={(e) => setTypingTestText(e.target.value)}
            required
            rows={6}
            className="w-full px-3 py-2 border rounded-md"
          />
        </div>
        <button
          type="submit"
          className="w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600"
        >
          Update Typing Test
        </button>
      </form>
      {message && (
        <div className="mt-4 p-2 bg-blue-100 text-blue-700 rounded">
          {message}
        </div>
      )}
    </div>
  );
};

export default UpdateTypingTestText;
