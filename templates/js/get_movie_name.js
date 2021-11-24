let fetchMovieName = async () => {
    const response = await fetch("/movie_name")
    const data = await response.json()
    return data.movie_name
}


async function getMovieName() {
    return await fetchMovieName()

}