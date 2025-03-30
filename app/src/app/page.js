"use client"; // Ensure this is a Client Component

import { useRouter } from "next/navigation";

export default function Home() {

    const router = useRouter();

    const goToLogin = () => {
        router.push("/login")
    };

    const goToRegister = () => {
        router.push("/register_user")
    };

    const goToSearch = () => {
        router.push("/search")
    };

    return (
        <div>
            <h1>TEIRx</h1>
            <button onClick={goToLogin}>Login</button>
            <button onClick={goToRegister}>Register</button>
            <button onClick={goToSearch}>Search</button>
        </div>
    );
}


