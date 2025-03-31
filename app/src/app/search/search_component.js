'use client'; // Use the client directive to use React hooks

import { useState } from 'react';
import { useRouter } from "next/navigation";

export default function SearchComponent({ initialText = "" }) {
    const [query, setQuery] = useState(initialText);
    const router = useRouter();

    const handleSearch = () => {
        if (query === "" || query.toLowerCase() === initialText.toLowerCase()) {
            return
        }
        router.push(`/search/${encodeURIComponent(query)}`)
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
        </div>
    )
}
