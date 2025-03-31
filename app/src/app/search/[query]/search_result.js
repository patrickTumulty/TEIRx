
'use client'; // Use the client directive to use React hooks

import { useEffect, useState } from "react"; // ✅ Import useEffect and useState
import { TEIRX_SERVER_URL, HttpStatus } from '@/core/global';
import Link from 'next/link';

function MovieItem({ data }) {
    return (
        <div>
            <Link href={`/films/${data.imdb_id}`}>{data.title}</Link>
            <img src={data.poster}></img>
        </div>
    );
}

export default function SearchPage({ query }) {

    const [searchResults, setSearchResults] = useState([])

    async function search() {

        console.log("Searching: " + query)

        const qString = `/search?query=${query}`;

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

    useEffect(() => {
        console.log("Search")
        search()
    }, []);

    return (
        <div>
            <ul className="search-results">
                {searchResults.map((item, index) => (
                    <li key={index} className="search-item">
                        <MovieItem data={item} />
                    </li>
                ))}
            </ul>
        </div>
    );
}
