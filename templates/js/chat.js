const input = document.querySelector('#namefield')
const messages = document.querySelector('#messages')
const username = document.querySelector('#namefield')

let ws = new WebSocket("ws://localhost:9000/chat");
ws.onmessage = function (msg) {
    insertMessage(JSON.parse(msg.data))
};


document.getElementById('chatButton').addEventListener('click', () => {
    const message = {
        typ: "add",
        userName: username.value,
    }
    ws.send(JSON.stringify(message));
    username.value = "";
    document.getElementById('chat').style.display = 'block'
    document.getElementById('name').style.display = 'none'
})

document.getElementById('sendButton').addEventListener('click', () => {
    let select = document.getElementById("inputState").value
    let result
    if (select == 'markdown') {
        result = marked.parse(textarea);
    }
    else {
        result = textarea
    }
    const message = {
        typ: "message",
        userName: username.value,
        msg: result
    }
    ws.send(JSON.stringify(message));
    username.value = "";
    document.getElementById("sendMessage").value = "";
    document.querySelector('textarea').style.cssText = 'height:65px'
})

function insertMessage(messageObj) {
    // Create a div object which will hold the message
    let getElementFromString = (string) => {
        let div = document.createElement('div');
        div.innerHTML = string;
        return div.firstElementChild;
    }
    let string
    if (messageObj.userName == username.value) {
        string = `<div class="container reciever mt-2 mb-2 border border-light bg-gradient   float-end  me-2" id="sender">
                  <h4 class="mt-2 ms-1 fw-bolder" style="font-family: sans-serif;">${messageObj.userName}</h4>
                  <p class=" ms-1 " style="font-family:sans-serif ; font-size: 15px;">${messageObj.msg}</p>
                  </div>`;
    }
    else {
        string = `<div class="container reciever mt-2 mb-2 border border-light bg-dark bg-gradient float-start  ms-2"
                   id="reciever">
                  <h4 class="mt-2 ms-1 fw-bolder" style="font-family: sans-serif;">${messageObj.userName}</h4>
                  <p class=" ms-1 " style="font-family:sans-serif ; font-size: 15px;">${messageObj.msg}</p>
                  </div>`;
    }

    // converting element string to dom
    let messageElement = getElementFromString(string);
    // console.log(parameterElement);
    messages.appendChild(messageElement);
}
