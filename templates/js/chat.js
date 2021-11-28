document.getElementById('chatButton').addEventListener('click', () => {
    document.getElementById('chatButton').disabled = true;
    let namefield = document.getElementById('namefield').value;
    function redirect() {
        window.location.href('/chats')
    }

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
                redirect()
            } else {
                alert("Error!!!!");
            }
        })

})