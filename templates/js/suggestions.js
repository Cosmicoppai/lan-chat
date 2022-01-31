document.getElementById('button').addEventListener('click', () => {
    document.getElementById('movieError').style.display = 'none'
    document.getElementById('dateError').style.display = 'none'
    document.getElementById('button').disabled = true;
    let moviefield = document.getElementById('moviefield').value;
    let datefield = document.getElementById('datefield').value;
    let descriptionfield = document.getElementById('descriptionfield').value;
    function clear() {
        document.getElementById("moviefield").value = "";
        document.getElementById("descriptionfield").value = "";
        document.getElementById("datefield").value = "";
        document.getElementById('button').disabled = false;
    }
    if (moviefield == '' && datefield == '') {
        document.getElementById('movieError').style.display = 'block'
        document.getElementById('dateError').style.display = 'clock'
        document.getElementById('button').disabled = false;
    }
    else if (moviefield == '') {
        document.getElementById('movieError').style.display = 'block'
        document.getElementById('button').disabled = false;
    }
    else if (datefield == '') {
        document.getElementById('dateError').style.display = 'block'
        document.getElementById('button').disabled = false;
    }
    else if (moviefield && datefield) {
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
                    }
                    else {
                        alert("Error!!!!");
                    }
               
            })
    }

})