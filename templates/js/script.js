let movie = localStorage.getItem("movie_name")
if (movie == null) {
    movieName()
    movie = localStorage.getItem("movie_name")
}
let movieSrc = document.getElementById('source')
movieSrc.setAttribute('src', '/get_movie/'+ movie);