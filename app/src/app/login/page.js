'use client'; // Use the client directive to use React hooks

import { useState } from 'react';
import { TEIRX_SERVER_URL, HttpStatus } from '@/core/global';

export default function Login() {
    const [token, setToken] = useState('')
    const [error, setError] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');

    const login = async (username, password) => {
        try {
            const res = await fetch(TEIRX_SERVER_URL + '/login', {
                method: 'POST',
                body: JSON.stringify({ 'username': username, 'password': password }),
                headers: {
                    'content-type': 'application/json'
                }
            })

            if (!res.ok) {
                setError("Incorrect username or password")
                return
            }

            const data = await res.json();
            setToken(data.token)
        } catch (error) {
            console.log(error)
        }
    }

    const handleLogin = (e) => {
        e.preventDefault();
        setError("")
        login(username, password)
    };


    const logout = async () => {
        const rsp = await fetch(TEIRX_SERVER_URL + '/logout', {
            method: 'POST',
            body: JSON.stringify({ 'token': token }),
            headers: {
                'content-type': 'application/json'
            }
        })
        if (rsp.ok) {
            setToken('')
        }
    }

    const handleLogout = (_) => {
        logout()
    }

    return (
        <div className='form-container'>
            <form className='basic-form' onSubmit={handleLogin} >
                <h2 >Login</h2>
                <label>Username</label>
                <input
                    type="text"
                    id="username"
                    name="username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder="Enter your username"
                />
                <label>Password</label>
                <input
                    type="password"
                    id="password"
                    name="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Enter your password"
                />
                {error !== '' && <p className="error-label">Login failed: {error}</p>}
                <button type="submit">Login</button>
            </form>
            {token !== '' && <button onClick={handleLogout}>Logout</button>}
        </div>
    );
}
