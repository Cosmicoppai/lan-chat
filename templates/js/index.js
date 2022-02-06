const toggleSwitch = document.querySelector('.theme-switch input[type="checkbox"]');

const switchTheme = (e) => {
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
const theme = localStorage.getItem('theme') ? localStorage.getItem('theme') : 'dark'
if (window?.matchMedia('(prefers-color-scheme: dark)').matches) {
    document.documentElement.setAttribute('data-theme', 'dark');
    if (theme === 'dark') {
        toggleSwitch.checked = true;
        document.getElementById('light').style.opacity = '0.5'
    }
}
else if (window?.matchMedia('(prefers-color-scheme: light)').matches) {
    document.documentElement.setAttribute('data-theme', 'light');
}

if (currentTheme) {
    document.documentElement.setAttribute('data-theme', currentTheme);
    if (currentTheme === 'dark') {
        toggleSwitch.checked = true;
        document.getElementById('light').style.opacity = '0.5'
    }
    else {
        document.getElementById('dark').style.opacity = '0.5'
    }
}

Object.defineProperty(String.prototype, 'capitalize', {
    value: function () {
        return this.charAt(0).toUpperCase() + this.slice(1);
    },
    enumerable: false
});


getMovieName().then(x => {
    let remove_after = x.indexOf('.');
    let result = x.substring(0, remove_after);
    document.getElementById('movieName').innerHTML += result.capitalize()
})
