'use client'; // Use the client directive to use React hooks

import { useState } from 'react';

export default function Login() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');


    const login = async (username, password) => {
        try {
            const res = await fetch('http://localhost:8080/login', {
                method: 'POST',
                body: JSON.stringify({ 'username': username, 'password': password }),
                headers: {
                    'content-type': 'application/json'
                }
            })
            console.log(res)
            if (res.ok) {
                console.log("Yeai!")
            } else {
                console.log("Oops! Something is wrong.")
            }
        } catch (error) {
            console.log(error)
        }
    }

    const handleLogin = (e) => {
        e.preventDefault();
        console.log('Username:', username);
        console.log('Password:', password);

        login(username, password)
    };

    return (
        <div className="min-h-screen bg-gray-100 flex items-center justify-center">
            <div className="bg-white p-8 rounded-lg shadow-lg w-full max-w-sm">
                <h2 className="text-2xl font-semibold text-center text-gray-700 mb-6">Login</h2>

                <form onSubmit={handleLogin} className="space-y-4">
                    <div>
                        <label htmlFor="username" className="block text-sm font-medium text-gray-600">Username</label>
                        <input
                            type="text"
                            id="username"
                            name="username"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            placeholder="Enter your username"
                            className="w-full mt-2 px-4 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
                        />
                    </div>

                    <div>
                        <label htmlFor="password" className="block text-sm font-medium text-gray-600">Password</label>
                        <input
                            type="password"
                            id="password"
                            name="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            placeholder="Enter your password"
                            className="w-full mt-2 px-4 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
                        />
                    </div>

                    <button
                        type="submit"
                        className="w-full py-2 bg-blue-500 text-white font-semibold rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-400"
                    >
                        Login
                    </button>
                </form>
            </div>
        </div>
    );
}
