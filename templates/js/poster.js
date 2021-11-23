let movie
if (localStorage.getItem("movie_name")) {
    movie = localStorage.getItem("movie_name")
}else {
    movieName()
    movie = localStorage.getItem("movie_name")
}
let imageSrc = document.getElementById('image')
imageSrc.setAttribute('src', '/get_poster/'+ movie);