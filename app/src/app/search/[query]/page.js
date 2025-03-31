
import SearchPage from './search_result'
import SearchComponent from '../search_component';

export default async function Search({ params }) {

    const { query } = await params

    return (
        <div>
            <SearchComponent initialText={decodeURIComponent(query)}/>
            <SearchPage query={query} />
        </div>
    );
}
