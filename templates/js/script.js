getMovieName().then(movie => {
    let movieSrc = document.getElementById('source')
    movieSrc.setAttribute('src', '/get_movie/' + movie);
})