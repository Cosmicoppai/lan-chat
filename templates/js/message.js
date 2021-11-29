document.getElementById('sendButton').addEventListener('click', () => {
    let textarea = document.getElementById('sendMessage').value;
    let select = document.getElementById("inputState").value
    let result
    if ( select == 'markdown'){
        result = marked.parse(textarea);
    }
    else{
        result = textarea
    }
    console.log(result)
    function clear() {
        document.getElementById("sendMessage").value = "";
        document.querySelector('textarea').style.cssText ='height:65px'
    }
    fetch("/",
        {
            method: "post",
            body: JSON.stringify({
                message:result
            }),
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => {
            if (response.ok) {
                clear()
            } else {
                alert("Couldn't send your message!");
                clear()
            }
        })

})