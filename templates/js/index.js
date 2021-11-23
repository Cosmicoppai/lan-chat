const toggleSwitch = document.querySelector('.theme-switch input[type="checkbox"]');

const switchTheme = (e) => {
    location.reload()
    if (e.target.checked) {
        document.documentElement.setAttribute('data-theme', 'dark');
        localStorage.setItem('theme', 'dark'); //add this
    }
    else {
        document.documentElement.setAttribute('data-theme', 'light');
        localStorage.setItem('theme', 'light'); //add this
    }
}
toggleSwitch.addEventListener('change', switchTheme, false);
const currentTheme = localStorage.getItem('theme') ? localStorage.getItem('theme') : null;

if (currentTheme) {
    document.documentElement.setAttribute('data-theme', currentTheme);

    if (currentTheme === 'dark') {
        toggleSwitch.checked = true;
        document.getElementById('light').style.opacity = '0.5'
        document.getElementById('blackButton').style.display='none';
        document.getElementById('whiteButton').style.display='block';
    }
    else {
        document.getElementById('dark').style.opacity = '0.5'
        document.getElementById('blackButton').style.display='block';
        document.getElementById('whiteButton').style.display='none';
    }
}


Object.defineProperty(String.prototype, 'capitalize', {
    value: function() {
        return this.charAt(0).toUpperCase() + this.slice(1);
    },
    enumerable: false
});

let x = localStorage.getItem("movie_name")
if (x === null) {
    movieName()
    x = localStorage.getItem("movie_name")
}

let remove_after = x.indexOf('.');
let result = x.substring(0, remove_after);
document.getElementById('movieName').innerHTML += result.capitalize()