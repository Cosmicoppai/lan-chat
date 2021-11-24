getMovieName().then(moviePoster => {
    let removeSuffix = moviePoster.indexOf('.');
    let Poster = moviePoster.substring(0, removeSuffix);
    let imageSrc = document.getElementById('poster')
    imageSrc.setAttribute('src', '/get_poster/' + Poster);
})