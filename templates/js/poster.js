let moviePoster = localStorage.getItem("movie_name")
if (moviePoster == null) { // if movie_name doesn't exist in localStorage
    movieName()
    moviePoster = localStorage.getItem("movie_name")
}
let removeSuffix = moviePoster.indexOf('.');
let Poster = moviePoster.substring(0, removeSuffix);
let imageSrc = document.getElementById('poster')
imageSrc.setAttribute('src', '/get_poster/'+ Poster);