const input = document.querySelector('#sendMessage')
const messages = document.querySelector('#messages')
const username = document.querySelector('#namefield')
const image = document.querySelector('#file')

let token

let usernameInput = document.getElementById("namefield");
let messageInput = document.getElementById("sendMessage");

usernameInput.addEventListener("keyup", function (event) {
    if (event.keyCode === 13) {
        event.preventDefault();
        document.getElementById("chatButton").click();
    }
});

messageInput.addEventListener("keyup", function (event) {
    if (event.keyCode === 13) {
        event.preventDefault();
        document.getElementById("sendButton").click();
    }
});


function startWebsocket() {
    let ws = new WebSocket(`ws://${document.domain}:9000`);
    ws.onmessage = function (msg) {
        insertMessage(JSON.parse(msg.data))
    };

    ws.onclose = function () {
        document.getElementById('messages').innerHTML = ""
        location.reload()
    };

    document.getElementById('file').addEventListener('change',() => {
        let data = document.getElementById("file").files[0];
        let getElementFromString = (string) => {
            let div = document.createElement('div');
            div.innerHTML = string;
            return div.firstElementChild;
        }
        let string = `<div>
                        <span class="me-2" id="fileName">Do you want to send ${data.name}</span>
                      </div>`
        let messageElement = getElementFromString(string);
        document.getElementById("dataName").appendChild(messageElement);
        setTimeout(() => {
            document.getElementById('fileAsk').style.display = 'block'
        }, 1000);
    }
    );

    document.getElementById('imgSend').addEventListener('click', ()=> {
        let file = image.files[0]
        let fsize = file.size;
        let size = Math.round((fsize / 1024));
        if (size < 5121) {
            let reader = new FileReader();
            reader.onload = function (e) {
                const message = {
                    typ: "img-msg",
                    msg: reader.result
                }
                ws.send(JSON.stringify(message));
            }
            reader.readAsDataURL(file);
            document.getElementById('fileAsk').style.display = 'none'
            document.getElementById('dataName').innerHTML = ""
        }
        // console.log(file)
    })

    document.getElementById('chatButton').addEventListener('click', () => {
        if (username.value !== '') {
            const message = {
                typ: "add",
                userName: username.value,
            }
            ws.send(JSON.stringify(message));
            document.getElementById('chat').style.display = 'block'
            document.getElementById('name').style.display = 'none'
            document.getElementById('homepageButton').style.display = 'none'
        }
    })
    document.getElementById('leaveButton').addEventListener('click', () => {
        const message = {
            typ: "remove",
            userName: username.value,
            token: token
        }
        ws.send(JSON.stringify(message));
        username.value = '';
        document.getElementById('messages').innerHTML = "";
        setTimeout(location.reload(), 1000)
    })
    document.getElementById('sendButton').addEventListener('click', () => {
        let select = document.getElementById("inputState").value
        let result
        if (input.value !== '') {
            if (select === 'markdown') {
                result = marked.parse(input.value);
            }
            else {
                result = input.value
            }
            const message = {
                typ: "txt-msg",
                msg: result
            }
            ws.send(JSON.stringify(message));
            document.getElementById("sendMessage").value = "";
            document.querySelector('textarea').style.cssText = 'height:65px'
        }
    })
    function insertMessage(messageObj) {

        setTimeout(() => {
            messages.scrollTop = messages.scrollHeight
        }, 0)
        let getElementFromString = (string) => {
            let div = document.createElement('div');
            div.innerHTML = string;
            return div.firstElementChild;
        }
        let string
        if (messageObj.typ === 'alert') {
            string = `<div class="mt-2 mb-2  float-start  ms-2 text-center" style="width: 95%;">
                          <span class="mt-2 ms-1 px-3 py-2 border border-light bg-dark rounded-pill"
                         style="font-family: sans-serif; font-size: 15px;"><span
                        style="font-family:sans-serif ; font-size: 15px;"><span class="fw-bolder">${messageObj.msg}</span>
                         </div>`
            token = messageObj.token !== "" ? messageObj.token : token
            document.getElementById("onlineUser").innerHTML = 'Online:' + messageObj.totalUser;
        }
        else if (messageObj.typ === 'txt-msg') {
            if (messageObj.userName === username.value) {
                string = `<div class="container reciever mt-2 mb-2 border border-light bg-gradient   float-end  me-2" id="sender"  >
                              <h4 class="mt-2 ms-1 fw-bolder" style="font-family: sans-serif;">${messageObj.userName}</h4>
                              <p class="ms-1 text-break" style="font-family:sans-serif ; font-size: 15px;">${messageObj.msg}</p>
                              </div>`;

            }
            else if ((messageObj.userName !== username.value)) {
                string = ` <div class="container  reciever mt-2 mb-2 border border-light bg-dark bg-gradient float-start  ms-2"
                id="reciever" >
                <h4 class="mt-2 ms-1 fw-bolder" style="font-family: sans-serif;">${messageObj.userName}</h4>
                <p class="ms-1 text-break" style="font-family:sans-serif ; font-size: 15px; ">${messageObj.msg}</p>
                </div>`;

            }
        }
        else if (messageObj.typ === 'img-msg') {
            if (messageObj.userName === username.value) {
                string = `  <div class="container reciever mt-2 mb-2 border border-light bg-gradient float-end me-2 " id="sender">
                            <h4 class="mt-2 ms-1 fw-bolder" style="font-family: sans-serif;">${messageObj.userName}</h4>
                            <p class=" mx-auto ms-1" ><img width="660" height="325" src=${messageObj.msg} onclick="imgClick()" id="imgSender" alt="" class="imgSender mx-auto"></p>
                            </div> `;
            }
            else if ((messageObj.userName !== username.value)) {
                string = `  <div class="container reciever mt-2 mb-2 border border-light bg-gradient float-start me-2 " id="sender">
                            <h4 class="mt-2 ms-1 fw-bolder" style="font-family: sans-serif;">${messageObj.userName}</h4>
                            <p class=" ms-1 mx-auto" ><img  width="660" height="325"src=${messageObj.msg} onclick="imgClick()" id="imgSender" alt="" class="imgSender mx-auto"></p>
                            </div> `;
            }


        }
        // converting element string to dom
        let messageElement = getElementFromString(string);
        // console.log(parameterElement);
        messages.appendChild(messageElement);
    }
}

startWebsocket();
