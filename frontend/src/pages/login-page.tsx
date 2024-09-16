import React, { useState } from 'react';
import { UserIcon, KeyIcon, LockOpenIcon, XCircle } from 'lucide-react';
import { server } from '@common/server';
import * as types from '@common/types';
import { Button } from '@/components/ui/button';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [userPassword, setUserPassword] = useState('');
  const [testPassword, setTestPassword] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    localStorage.setItem('username', username);
    localStorage.setItem('userPassword', userPassword);
    localStorage.setItem('testPassword', testPassword);

    server.send_message({
      Typ: types.Varient.UserLogin,
      Val: {
        Username: username,
        Password: userPassword,
        TestCode: testPassword,
      }
    });
    console.log('Login submitted:', { username, userPassword, testPassword });
  };

  const handleQuit = () => {
    server.send_message({
      Typ: types.Varient.Quit,
      Val: {}
    });
  };

  return (
    <div className="min-h-screen bg-blue-700 flex flex-col lg:flex-row items-center justify-around p-4">

      <div className="text-white mb-8 lg:mb-0 lg:mr-8 text-center lg:text-left">
        <img src="/WCL_LOGO.png" alt="Coal India Logo" className="w-24 h-24 mx-auto lg:mx-0 mb-4" />
        <h1 className="text-3xl font-bold mb-2">Welcome to</h1>
        <h2 className="text-2xl font-semibold mb-2">Western Coalfields Limited (WCL)</h2>
        <h3 className="text-xl">Computer Aptitude Test</h3>
      </div>


      <div className="bg-white rounded-lg shadow-xl p-8 w-full max-w-md">
        <h2 className="text-2xl font-bold text-blue-700 mb-6 text-center">Login</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="username" className="block text-sm font-medium text-gray-700 mb-1">Username</label>
            <div className="relative">
              <UserIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={20} />
              <input
                id="username"
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="pl-10 w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter your username"
                required
              />
            </div>
          </div>
          <div>
            <label htmlFor="userPassword" className="block text-sm font-medium text-gray-700 mb-1">User Password</label>
            <div className="relative">
              <KeyIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={20} />
              <input
                id="userPassword"
                type="password"
                value={userPassword}
                onChange={(e) => setUserPassword(e.target.value)}
                className="pl-10 w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter your password"
                required
              />
            </div>
          </div>
          <div>
            <label htmlFor="testPassword" className="block text-sm font-medium text-gray-700 mb-1">Test Password</label>
            <div className="relative">
              <LockOpenIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={20} />
              <input
                id="testPassword"
                type="password"
                value={testPassword}
                onChange={(e) => setTestPassword(e.target.value)}
                className="pl-10 w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter test password"
                required
              />
            </div>
          </div>
          <div className="flex justify-between">
            <button
              type="submit"
              className="bg-green-600 text-white py-2 px-4 rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 transition duration-150 ease-in-out"
            >
              LOGIN
            </button>
            <Button
              type="button"
              onClick={handleQuit}
              variant="destructive"
              className="flex items-center"
            >
              <XCircle className="mr-2" size={16} />
              Quit
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
