let getElementFromString = (string) => {
    let li = document.createElement('li');
    li.innerHTML = string;
    return li.firstElementChild;
}

getDetails().then(movie => {
    if (movie === null) {
        document.getElementById('normal').style.display = 'none'
        document.getElementById('main2').style.display = 'table'
    }
    else {
        for (i = 0; i <= 2; i++) {
            if (movie[i]) {
                let link
                let imglink = movie[i].imageLink.split(' ').join('%20')
                let name = movie[i].name.split(' ').join('+')
                if (movie[i].typ === 'movie') {
                    link = `/movie?movieName=${name}`
                }
                else if (movie[i].typ === 'series') {
                    link = `/movie?seriesName=${name}&ep${movie[i].epNo}`
                }
                document.getElementById(`img${i}`).setAttribute('src', imglink);

                document.getElementById(`name${i}`).innerHTML += movie[i].epNo?`${movie[i].name}(Ep${movie[i].epNo})`:`${movie[i].name}`
            }
        }
        for (i = 3; i <= movie.length; i++) {
            if (movie[i]) {
                let link
                let imglink = movie[i].imageLink.split(' ').join('%20')
                let name = movie[i].name.split(' ').join('+')
                if (movie[i].typ === 'movie') {
                    link = `/movie?movieName=${name}`
                }
                else if (movie[i].typ === 'series') {
                    link = `/movie?seriesName=${name}&ep=${movie[i].epNo}`
                }
                let header = document.getElementById('searchresult');
                let string = `<li>
                                    <a href=${link} title=${movie[i].name}  style="text-decoration: none;">
                                        <div class="film-poster" id="poster" >
                                            <span class="text">
                                                <img class="playImg" src="static/images/play1.png" alt="">
                                            </span>
                                            <img 
                                                class="resultimg" alt=${movie[i].name}
                                                src=${imglink}>
                                                ${movie[i].epNo ? ` < div className = "bottom-left" > Ep ${movie[i].epNo} < /div>`: ``}
                                        </div>
                                    </a>
                                    <div class="details">
                                        <a href=${link} title="${movie[i].name}"  style="text-decoration: none;">
                                             <p class="name">${movie[i].name}</p>
                                        </a>
                                    </div>
                                  </li>`;
                let headerElement = getElementFromString(string);
                header.appendChild(headerElement);
            }
        }
    }

})

