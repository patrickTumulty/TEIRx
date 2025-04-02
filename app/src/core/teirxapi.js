import { TEIRX_SERVER_URL, HttpStatus } from '@/core/global';

export class TeirxApi {

    static #getUrl(str) {
        return `${TEIRX_SERVER_URL}${str}`
    }

    static async search(query) {
        try {
            const res = await fetch(this.#getUrl(`/search?query=${query}`), {
                method: 'GET',
                headers: {
                    'content-type': 'application/json'
                }
            })

            if (!res.ok) {
                console.log("Search error")
                return
            }

            return await res.json()
        } catch (error) {
            console.log("Error searching for film: " + error)
        }
    }

    static async getFilmData(filmID) {
        try {
            const res = await fetch(this.#getUrl(`/get-film?id=${filmID}`), {
                method: 'GET',
                headers: {
                    'content-type': 'application/json'
                }
            })

            if (!res.ok) {
                console.log("Get film error")
                return
            }

            return await res.json()
        } catch (error) {
            console.log("Error getting film data: " + error)
        }
    }
};
