let getElementFromString = (string) => {
    let li = document.createElement('li');
    li.innerHTML = string;
    return li.firstElementChild;
}

getDetails().then(movie => {
    function appendElements() {
        for (i = 0; i <= 3; i++) {
            if (movie[i]) {
                let link
                if(movie[i].typ ==='movie'){
                    link = `/movie?movieName=${movie[i].name}`
                }
                else if(movie[i].typ ==='series'){
                    link = `/movie?seriesName=${movie[i].name}&ep${movie[i].epNo}`
                }
                let header = document.getElementById('slideWrap');
                let string = `<li><a href=${link} class="column caption  col-xs-6" id="caption"><span class="text">
                <img class="playImg" src="static/images/play1.png" alt="">
            </span><img src=${movie[i].imageLink} id='img'></a></li>`;
                let headerElement = getElementFromString(string);
                header.appendChild(headerElement);
            }
        }
        for (i = 0; i <= movie.length; i++) {
            if (movie[i]) {
                let link
                if(movie[i].typ ==='movie'){
                    link = `/movie?movieName=${movie[i].name}`
                }
                else if(movie[i].typ ==='series'){
                    link = `/movie?seriesName=${movie[i].name}&ep=${movie[i].epNo}`
                }
                let header = document.getElementById('searchresult');
                let string = `<li>
                                <a href=${link} title="${movie[i].name}"  style="text-decoration: none;">
                                    <div class="film-poster" id="poster" >
                                        <span class="text">
                                            <img class="playImg" src="static/images/play1.png" alt="">
                                        </span>
                                        <img 
                                            class="resultimg" alt="Yami Shibai 10"
                                            src=${movie[i].imageLink}>
                                        <div class="bottom-left">Ep ${movie[i].epNo}</div>
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
    async function getDetails() {
        await appendElements()
        while (document.getElementById('img').src!==null) {
            let responsiveSlider = function () {
                let slider = document.getElementById("slider");
                let sliderWidth = slider.offsetWidth;
                let slideList = document.getElementById("slideWrap");
                let count = 1;
                let items = slideList.querySelectorAll("li").length;
                let prev = document.getElementById("prev");
                let next = document.getElementById("next");
                window.addEventListener('resize', function () {
                    sliderWidth = slider.offsetWidth;
                });
                let prevSlide = function () {
                    if (count > 1) {
                        count = count - 2;
                        slideList.style.left = "-" + count * sliderWidth + "px";
                        count++;
                    }
                    else if (count = 1) {
                        count = items - 1;
                        slideList.style.left = "-" + count * sliderWidth + "px";
                        count++;
                    }
                };
                let nextSlide = function () {
                    if (count < items) {
                        slideList.style.left = "-" + count * sliderWidth + "px";
                        count++;
                    }
                    else if (count = items) {
                        slideList.style.left = "0px";
                        count = 1;
                    }
                };
                next.addEventListener("click", function () {
                    nextSlide();
                });
                prev.addEventListener("click", function () {
                    prevSlide();
                });
    
                let theInterval
                function startSlide() {
                    theInterval = setInterval(nextSlide(), 7000);
                }
                function stopSlide() {
                    clearInterval(theInterval);
                }
                $(function () {
                    startSlide();
                    $('#slider').hover(function () {
                        stopSlide();
                    }, function () {
                        startSlide();
                    })
                });
            };
            await responsiveSlider()
            break;
        }

    }
    getDetails()
})

