"use client"; 

import { TeirxApi } from "@/core/teirxapi";
import { useState, useEffect } from 'react';
export default function Home() {

    const [filmData, setFilmData] = useState({})
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        TeirxApi.getFeatured().then((response) => {
            setFilmData(response)
            setLoading(false)
        })
    }, [])

    if (loading) {
        return <div>Loading... </div>
    }

    return (
        <div className="featured-film">
            <h1>{filmData.title} ({filmData.year})</h1>
            <img src={filmData.poster} />
        </div>
    );
}


