import { UserIcon, KeyIcon,  } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useStateContext } from '@/context/app-context';
import { server } from '@common/server';
import * as types from '@common/types';



export default function LoginPage() {

  const {
    username, 
    userPassword, 
    setUsername, 
    setUserPassword,
  } = useStateContext();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    server.send_message({
      Typ: types.Varient.UserLoginRequest,
      Val: {
        Username: username,
        Password: userPassword,
      }
    });
    console.log('Login submitted:', { username, userPassword });
  };

  const handleQuit = () => {
    server.send_message({
      Typ: types.Varient.Quit,
      Val: {}
    });
  };

  const handleCheckSystem = () => {
    server.send_message({
      Typ: types.Varient.CheckApps,
      Val: {}
    });
  }
  
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
          <div className="flex justify-between">
            <Button
              type="submit"
              className="bg-green-600 text-white py-2 px-4 rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 transition duration-150 ease-in-out"
            >
              Login
            </Button>
            <Button
              type="button"
              variant="default"
              onClick={handleCheckSystem}
            >
              Check System
            </Button>
            <Button
              type="button"
              onClick={handleQuit}
              variant="destructive"
              className="flex items-center"
            >
              Quit
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
