'use client'; // Use the client directive to use React hooks

import './search_component.css'
import { useState, useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';

export default function SearchComponent({ initialText = "" }) {
    const pathname = usePathname()
    const [query, setQuery] = useState(initialText);
    const router = useRouter();

    useEffect(() => {
        const match = pathname.match(/^\/search\/(.+)/);
        if (match && match[1]) {
            const decoded = decodeURIComponent(match[1]);
            setQuery(decoded);
        } else {
            setQuery('');
        }
    }, [pathname]);

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
