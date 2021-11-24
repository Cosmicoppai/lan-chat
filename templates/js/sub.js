getMovieName().then(subFileName => {
    let remove_after = subFileName.indexOf('.');
    let result = subFileName.substring(0, remove_after);

    const subtitles = (result) => {
        let subtitlesSrc = document.getElementById('subtitle')
        subtitlesSrc.setAttribute('src', '/get_sub/' + result);
    }
    subtitles(result)
})


