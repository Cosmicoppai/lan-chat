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

