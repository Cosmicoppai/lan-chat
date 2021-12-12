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

const fileFunction = ()=>{
        let data = document.getElementById("file").files[0];
        document.getElementById('fileName').innerHTML = 'Do you want to send '+ data.name
        setTimeout(() => {
            document.getElementById('fileAsk').style.display='block'
        }, 100);
}

const fileSend =()=>{
    document.getElementById('fileAsk').style.display='none'
    document.getElementById('fileName').value = ''
    let data = document.getElementById("file").value =''
}


let modal = document.getElementById("myModal");
let img = document.getElementById("imgSender");
let img2 = document.getElementById("imgSender2");
let modalImg = document.getElementById("img01");
let captionText = document.getElementById("caption");
img.onclick = function(){
  modal.style.display = "block";
  modalImg.src = this.src;
  captionText.innerHTML = this.alt;
}
img2.onclick = function(){
  modal.style.display = "block";
  modalImg.src = this.src;
  captionText.innerHTML = this.alt;
}
let span = document.getElementsByClassName("close")[0];
span.onclick = function() {
  modal.style.display = "none";
}
