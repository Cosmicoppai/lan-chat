document.getElementById('chatButton').addEventListener('click', () => {
    document.getElementById('chatButton').disabled = true;
    let namefield = document.getElementById('namefield').value;
   

    fetch("/",
        {
            method: "post",
            body: JSON.stringify({
                name: namefield
            }),
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => {
            if (response.ok) {
                document.getElementById('chatButton').disabled = false;
                document.getElementById('chat').style.display = 'block'
                document.getElementById('name').style.display = 'none'
                
            } else {
                alert("Error!!!!");
            }
        })

})
