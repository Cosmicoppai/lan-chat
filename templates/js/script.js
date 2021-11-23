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
        movieSrc.setAttribute('src', '/get_movie/' + x);
        fetch('/get_sub/'+ result)
            .then(data => {
                return data.json();
            })
            .then(post => {
                let x = post['subtitle']
                let subtitlesSrc = document.getElementById('subtitle')
                subtitlesSrc.setAttribute('src', '/get_sub/' + x);
            });
    });


