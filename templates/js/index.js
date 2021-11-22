const toggleSwitch = document.querySelector('.theme-switch input[type="checkbox"]');

const switchTheme= (e)=> {
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
        document.getElementById('light').style.opacity ='0.5'
    }
    else{
        document.getElementById('dark').style.opacity ='0.5'
        
    }
}

var today = new Date();
var dd = today.getDate();
var mm = today.getMonth()+1; //January is 0 so need to add 1 to make it 1!
var yyyy = today.getFullYear();
if(dd<10){
  dd='0'+dd
} 
if(mm<10){
  mm='0'+mm
} 

today = yyyy+'-'+mm+'-'+dd;
document.getElementById("datefield").setAttribute("min", today);

fetch('/movie_name')
.then(data => {
return data.json();
})
.then(post => {
console.log(post.movie_name);
document.getElementById('movieName').innerHTML += post.movie_name
let movieSrc = document.getElementById('source')
movieSrc.setAttribute('src', '/movie/'+ post.movie_name);
});