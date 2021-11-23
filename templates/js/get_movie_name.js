let movieName = () => {
    let x
    fetch('/movie_name')
        .then(data => {
            return data.json();
        })
        .then(post => {
            x = post.movie_name
            localStorage.setItem("movie_name", x)
        });

}


function getMovieName() {
    let movie = localStorage.getItem("movie_name")
    if (movie === null) {
        movieName()
        movie = localStorage.getItem("movie_name")

    }
    return movie
}