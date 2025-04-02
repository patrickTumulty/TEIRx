import FilmInfo from "./film_info";


export default async function FilmPage({ params }) {

    const { film } = await params

    return (
        <div>
            <FilmInfo filmID={film}/>
        </div>
    );
}
