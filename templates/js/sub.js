let subFileName = localStorage.getItem("movie_name") // get the movie_name from localStorage
let remove_after = subFileName.indexOf('.');
let result = subFileName.substring(0, remove_after);

const subtitles = (result)=> {
    let subtitlesSrc = document.getElementById('subtitle')
    subtitlesSrc.setAttribute('src', '/get_sub/' + result);
}
subtitles(result)


