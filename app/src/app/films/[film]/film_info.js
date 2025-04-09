
'use client';

import { TeirxApi } from "@/core/teirxapi";
import react from "react";

export default function FilmInfo({ filmID }) {

    const [filmData, setFilmData] = react.useState({})
    const [loading, setLoading] = react.useState(true)
    const ranks = ['S', 'A', 'B', 'C', 'D', 'F']

    react.useEffect(() => {
        console.log("Getting film data")
        TeirxApi.getFilmData(filmID).then((response) => {
            setFilmData(response)
            setLoading(false)
        })
    }, []);

    if (loading) {
        return <div>Loading... </div>
    }

    return (
        <div>
            <img src={filmData.poster} />
            <h1>{filmData.title} ({filmData.year})</h1>
            <p>{filmData.plot}</p>
            <h2>Rank: {filmData.stats.rank}</h2>
            <h3>Total: {filmData.stats.total_count}</h3>
            <h4>[S: {filmData.stats.s_count}, A: {filmData.stats.a_count}, B: {filmData.stats.b_count}, C: {filmData.stats.c_count}, D: {filmData.stats.d_count}, F: {filmData.stats.f_count}]</h4>
            <ul>
                {ranks.map((item, index) => (
                    <li key={index}>
                        <button>{item}</button>
                    </li>
                ))}
            </ul>
        </div>

    );
}


