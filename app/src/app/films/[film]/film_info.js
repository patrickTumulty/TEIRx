
'use client';

import { TeirxApi } from "@/core/teirxapi";
import react from "react";

export default function FilmInfo({ filmID }) {

    const [filmData, setFilmData] = react.useState({})
    const [loading, setLoading] = react.useState(true)
    
    react.useEffect(() => {
        console.log("Getting film data")
        TeirxApi.getFilmData(filmID).then((response) => {
            console.log("Getting: " + JSON.stringify(response, null, 2))
            setFilmData(response)
            setLoading(false)
        })
    }, []);

    if (loading) {
        return <div>Loading... </div>
    }

    return (
        <div>
            <img src={filmData.poster}/>
            <h1>{filmData.title} ({filmData.year})</h1>
            <p>{filmData.plot}</p>
            <h2>Rank: {filmData.stats.rank}</h2>
            <h3>Total: {filmData.stats.total_count}</h3>
            <h3>S: {filmData.stats.s_count}</h3>
            <h3>A: {filmData.stats.a_count}</h3>
            <h3>B: {filmData.stats.b_count}</h3>
            <h3>C: {filmData.stats.c_count}</h3>
            <h3>D: {filmData.stats.d_count}</h3>
            <h3>F: {filmData.stats.f_count}</h3>
        </div>
    );
}


