document.getElementById('button').addEventListener('click', () => {
    let moviefield = document.getElementById('moviefield').value;
    let datefield = document.getElementById('datefield').value;
    let descriptionfield = document.getElementById('descriptionfield').value;
    function clear() {
        document.getElementById("moviefield").value = "";
        document.getElementById("descriptionfield").value = "";
        document.getElementById("datefield").value = "";
    }
    console.log(moviefield, datefield, descriptionfield)
    let formData = new URLSearchParams();
    formData.append('movie_name', moviefield);
    formData.append('date', datefield);
    formData.append('msg', descriptionfield);
    fetch("/send-suggestion",
    {   
        body: formData,
        method: "post",
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }).then(response => {
        if (response.ok) { 
            alert("Suggestion sent successfully!")
            clear()
          } else {
            alert("Error!!!!");
          }
    })
    
})