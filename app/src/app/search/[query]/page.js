
import SearchPage from './search_result'

export default async function Search({ params }) {

    const { query } = await params

    return (
        <div>
            <SearchPage query={query} />
        </div>
    );
}
