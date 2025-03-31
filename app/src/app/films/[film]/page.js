
export default async function Page({ params }) {

    const { film } = await params

    return (
        <div>
            <h1>Movie: {film}</h1>
            <p>Details about {film}...</p>
        </div>
    );
}
