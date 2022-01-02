const input = document.querySelector('#sendMessage')
const messages = document.querySelector('#messages')
const username = document.querySelector('#namefield')
let token

let ws = new WebSocket("ws://localhost:9000/chat");

ws.onmessage = function (msg) {
    insertMessage(JSON.parse(msg.data))
};

document.getElementById('chatButton').addEventListener('click', () => {
    if (username.value != '') {
        const message = {
            typ: "add",
            userName: username.value,
        }
        ws.send(JSON.stringify(message));
        document.getElementById('chat').style.display = 'block'
        document.getElementById('name').style.display = 'none'
    }
})


document.getElementById('leaveButton').addEventListener('click', () => {
    console.log(token)
    const message = {
        typ: "remove",
        userName: username.value,
        token: token
    }
    ws.send(JSON.stringify(message));
    username.value ='';
    document.getElementById('messages').innerHTML = "";
    document.getElementById('chat').style.display = 'none'
    document.getElementById('name').style.display = 'block'
})

document.getElementById('sendButton').addEventListener('click', () => {
    let select = document.getElementById("inputState").value
    let result
    if (input.value != '') {
        if (select == 'markdown') {
            result = marked.parse(input.value);
        }
        else {
            result = input.value
        }
        const message = {
            typ: "message",
            msg: result
        }
        ws.send(JSON.stringify(message));
        document.getElementById("sendMessage").value = "";
        document.querySelector('textarea').style.cssText = 'height:65px'
    }
})

function insertMessage(messageObj) {
    console.log(messageObj)
    // Create a div object which will hold the message
    let getElementFromString = (string) => {
        let div = document.createElement('div');
        div.innerHTML = string;
        return div.firstElementChild;
    }
    let string
    if (messageObj.typ == 'alert') {
        string = `<div class="mt-2 mb-2  float-start  ms-2 text-center" style="width: 95%;">
                      <span class="mt-2 ms-1 px-3 py-2 border border-light bg-dark rounded-pill"
                     style="font-family: sans-serif; font-size: 15px;"><span
                    style="font-family:sans-serif ; font-size: 15px;"><span class="fw-bolder">${messageObj.msg}</span>
                     </div>`
        token = messageObj.token
        document.getElementById("onlineUser").innerHTML = 'Online:' + messageObj.totalUser;
    }
    else if (messageObj.typ == 'message'){
        if (messageObj.userName == username.value) {
            string = `<div class="container reciever mt-2 mb-2 border border-light bg-gradient   float-end  me-2" id="sender">
                          <h4 class="mt-2 ms-1 fw-bolder" style="font-family: sans-serif;">${messageObj.userName}</h4>
                          <p class=" ms-1 " style="font-family:sans-serif ; font-size: 15px;">${messageObj.msg}</p>
                          </div>`;
            
        }
        else if ((messageObj.userName != username.value)) {
            string = `<div class="container reciever mt-2 mb-2 border border-light bg-dark bg-gradient float-start  ms-2"
                           id="reciever">
                          <h4 class="mt-2 ms-1 fw-bolder" style="font-family: sans-serif;">${messageObj.userName}</h4>
                          <p class=" ms-1 " style="font-family:sans-serif ; font-size: 15px;">${messageObj.msg}</p>
                          </div>`;
    
        }
    }


    // converting element string to dom
    let messageElement = getElementFromString(string);
    // console.log(parameterElement);
    messages.appendChild(messageElement);
}
