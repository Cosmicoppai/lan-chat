let movie
if (localStorage.getItem("movie_name")) {
    movie = localStorage.getItem("movie_name")
}else {
    movieName()
    movie = localStorage.getItem("movie_name")
}
let movieSrc = document.getElementById('source')
movieSrc.setAttribute('src', '/get_movie/'+ movie);