import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { UserIcon, KeyIcon } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useToast } from "@/hooks/use-toast"
import { useNavigate } from 'react-router-dom';


export default function Login() {
    const [username, setUsername] = useState('');
    const [userPassword, setUserPassword] = useState('');
    const { toast } = useToast()
    const navigate = useNavigate();

    useEffect(() => {
        const checkAuthStatus = async () => {
            try {
                const response = await axios.get(`${import.meta.env.SERVER_URL}/admin/auth-status`, {
                    withCredentials: true,
                });
                
                if(response.data.isAuthenticated){
                    navigate("/dashboard");
                }
            } catch (err) {
                console.log(err);
            }
        }
        checkAuthStatus();
    }, [])

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            const response = await axios.post(`${import.meta.env.SERVER_URL}/admin/login`, {
                username,
                password: userPassword
            }, {
                withCredentials: true,
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            console.log('Login successful:', response.data);
            toast({
                title: "Login Successful",
                description: "You have been logged in successfully.",
            })

            navigate('/dashboard');


        } catch (error) {
            console.error('Login failed:', error);
            toast({
                variant: "destructive",
                title: "Login Failed",
                description: "Please check your credentials and try again.",
            })
        }
    };

    return (
        <div className="min-h-screen bg-blue-700 flex flex-col lg:flex-row items-center justify-around p-4">
            <div className="text-white mb-8 lg:mb-0 lg:mr-8 text-center lg:text-left">
                <img src="/WCL_LOGO.png" alt="Coal India Logo" className="w-24 h-24 mx-auto lg:mx-0 mb-4" />
                <h1 className="text-3xl font-bold mb-2">Welcome to</h1>
                <h2 className="text-2xl font-semibold mb-2">Western Coalfields Limited (WCL)</h2>
                <h3 className="text-xl">Computer Aptitude Test (Admin Pannel)</h3>
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
                        <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">Password</label>
                        <div className="relative">
                            <KeyIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={20} />
                            <input
                                id="password"
                                type="password"
                                value={userPassword}
                                onChange={(e) => setUserPassword(e.target.value)}
                                className="pl-10 w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                                placeholder="Enter your password"
                                required
                            />
                        </div>
                    </div>
                    <div className="flex w-full">
                        <Button
                            type="submit"
                            className="bg-green-600 text-white py-2 px-4 rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 transition duration-150 ease-in-out w-full"
                        >
                            LOGIN
                        </Button>
                    </div>
                </form>
            </div>
        </div>
    );
}
