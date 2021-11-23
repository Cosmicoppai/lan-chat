fetch('/movie_name')
    .then(data => {
        return data.json();
    })
    .then(post => {
        let x = post['movie-name']
        let remove_after = x.indexOf('.');
        let result = x.substring(0, remove_after);
        document.getElementById('movieName').innerHTML += result
        let movieSrc = document.getElementById('source')
        movieSrc.setAttribute('src', '/movie/' + result);
    });