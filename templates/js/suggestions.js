document.getElementById('button').addEventListener('click', () => {
    let moviefield = document.getElementById('moviefield').value;
    let datefield = document.getElementById('datefield').value;
    let descriptionfield = document.getElementById('descriptionfield').value;
    function redirect() {
        window.location.href = "/";
    }
    console.log(moviefield, datefield, descriptionfield)
    let formData = new FormData();
    formData.append('movie_name', moviefield);
    formData.append('date', datefield);
    formData.append('msg', descriptionfield);
    fetch("/send-suggestions",
    {   
        body: formData,
        method: "post",
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }).then(response => {
        if (response.ok) { 
            redirect()
          } else {
            alert("Error!!!!");
          }
    })
    
})