'use client'; // Use the client directive to use React hooks

import { TEIRX_SERVER_URL } from '@/core/global';
import { useState } from 'react';

export default function RegisterUser() {
    const [firstname, setFirstname] = useState('');
    const [lastname, setLastname] = useState('');
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const registerUser = async () => {
        try {
            const res = await fetch(TEIRX_SERVER_URL + "/register-user", {
                method: 'POST',
                body: JSON.stringify({
                    'username': username,
                    'firstname': firstname,
                    'lastname': lastname,
                    'email': email,
                    'password': password
                }),
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

    const handleRegisterUser = (e) => {
        e.preventDefault();
        registerUser(username, password)
    };

    return (
        <div className="form-container">
            <form className="basic-form" onSubmit={handleRegisterUser}>
                <label>Username</label>
                <input
                    type="text"
                    id="username"
                    name="username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder="Username"
                />
                <label>Email</label>
                <input
                    type="text"
                    id="email"
                    name="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="Email"
                />
                <div className='side-by-side'>
                    <label className='l' >First Name
                        <input
                            type="text"
                            id="firstname"
                            name="firstname"
                            value={firstname}
                            onChange={(e) => setFirstname(e.target.value)}
                            placeholder="First Name"
                        />
                    </label>
                    <label >Last Name
                        <input
                            type="text"
                            id="lastname"
                            name="lastname"
                            value={lastname}
                            onChange={(e) => setLastname(e.target.value)}
                            placeholder="Last Name"
                        />
                    </label>
                </div>
                <label>Password</label>
                <input
                    type="password"
                    id="password"
                    name="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Password"
                />
                <button type="submit">
                    Submit
                </button>
            </form>
        </div>
    );
}
