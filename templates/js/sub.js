getDetails().then(subFileName => {
    let url = window.location.href;
    url = url.split('/movie').pop();
    const argArray = new URLSearchParams(url)
    let movieName = argArray.has('movieName') ? argArray.get('movieName') : '';
    let seriesName = argArray.has('seriesName') ? argArray.get('seriesName') : '';
    let ep = argArray.has('ep') ? argArray.get('ep') : '';
    if (movieName !== undefined || seriesName !== undefined) {
        if (movieName) {
            const obj = subFileName.find(o => o.name === movieName)
            if (obj) {
                let movieSrc = document.getElementById('source')
                movieSrc.setAttribute('src', obj.videoLink);
                let subtitlesSrc = document.getElementById('subtitle')
                subtitlesSrc.setAttribute('src', obj.subLink);
            }
            else {
                document.getElementById('normal').style.display = 'none'
                document.getElementById('main').style.display = 'table'
            }
        }
        else if (seriesName && ep) {
            const obj = subFileName.find(o =>  o.name === seriesName && o.epNo === parseInt(ep) )
            if (obj) {
                let seriesSrc = document.getElementById('source')
                seriesSrc.setAttribute('src', obj.videoLink);
                let subtitlesSrc = document.getElementById('subtitle')
                subtitlesSrc.setAttribute('src', obj.subLink);
            }
            else {
                document.getElementById('normal').style.display = 'none'
                document.getElementById('main').style.display = 'table'
            }
        }
        else {
            document.getElementById('normal').style.display = 'none'
            document.getElementById('main').style.display = 'table'
        }
    }
    else {
        document.getElementById('normal').style.display = 'none'
        document.getElementById('main').style.display = 'table'
    }
})


