
'use client';

import react from "react";
import Link from 'next/link';
import { TeirxApi } from '@/core/teirxapi';

function MovieItem({ data }) {
    return (
        <div>
            <Link href={`/films/${data.imdb_id}`}>{data.title}</Link>
            <img src={data.poster}></img>
        </div>
    );
}

export default function SearchPage({ query }) {

    const [searchResults, setSearchResults] = react.useState([])
    const [loading, setLoading] = react.useState(true)

    react.useEffect(() => {
        TeirxApi.search(query).then((data) => {
            setSearchResults(data)
            setLoading(false)
        })
    }, []);

    if (loading) {
        return <div>Loading... </div>
    }

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
