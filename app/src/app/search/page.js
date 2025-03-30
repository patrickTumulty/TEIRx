'use client'; // Use the client directive to use React hooks

import { useState } from 'react';
import { TEIRX_SERVER_URL, HttpStatus } from '@/core/global';


function MovieItem({ data }) {
    return (
        <div>
            <label>{data.title}</label>
            <img src={data.poster}></img>
        </div>
    );
}


export default function Search() {

    const [query, setQuery] = useState('');
    const [searchResults, setSearchResults] = useState([])

    const search = async () => {

        const qString = `/search?query=${encodeURIComponent(query)}`;

        console.log(qString)

        try {
            const res = await fetch(TEIRX_SERVER_URL + qString, {
                method: 'GET',
                headers: {
                    'content-type': 'application/json'
                }
            })

            if (!res.ok) {
                console.log("Response error")
                return
            }

            const data = await res.json();
            setSearchResults(data)
        } catch (error) {
            console.log(error)
        }
    }

    const handleSearch = () => {
        search()
    };

    return (
        <div className="search-container">
            <input
                type="text"
                placeholder="Search..."
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                className="search-input"
            />
            <button onClick={handleSearch} className="search-button">
                Search
            </button>
            <ul className="search-results">
                {searchResults.map((item, index) => (
                    <li key={index} className="search-item">
                        <MovieItem data={item}/>
                    </li>
                ))}
            </ul>
        </div>
    );
}
