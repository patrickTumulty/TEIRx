import "./nav.css";
import SearchComponent from "./search/search_component";

export default function Navbar() {

    return (
        <nav className="nav">
            <a href="/" className="site-title">Redpen</a>
            <ul>
                <li>
                    <a href="/login">Login</a>
                </li>
                <li>
                    <a href="/register_user">New User</a>
                </li>
                <li>
                    <SearchComponent />
                </li>
            </ul>
        </nav>
    );
}
