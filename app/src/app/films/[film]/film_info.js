
'use client';

import { TeirxApi } from "@/core/teirxapi";
import react from "react";

export default function FilmInfo({ filmID }) {

    const [filmData, setFilmData] = react.useState({})
    react.useEffect(() => {
        TeirxApi.getFilmData(filmID).then((response) => {
            setFilmData(response)
        })
    }, [])

    return (
        <div>
            <h1>Movie: {filmData.title} ({filmData.year})</h1>
            <p>{filmData.plot}</p>
        </div>
    );
}
