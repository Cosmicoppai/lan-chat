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



document.getElementById('chatButton').addEventListener('click', () => {
    if (username.value !== '') {
        let ws = new WebSocket(`ws://${document.domain}:9000`)
        let data
        ws.onopen = function () {
            addUser()
        }
        ws.onmessage = function (msg) {
            data = JSON.parse(msg.data)
            if (data.typ === 'error') {
                alertMsg(JSON.parse(msg.data))
            }
            else {
                insertMessage(JSON.parse(msg.data))
                document.getElementById('chat').style.display = 'block'
                document.getElementById('name').style.display = 'none'
                document.getElementById('homepageButton').style.display = 'none'
            }
        };


        ws.onclose = function () {
            closeMsg()
        };

        function closeMsg() {
            if (!window.alert('Sorry some error occured!!!')) { window.location.reload(); }
        }
        function addUser() {
            if (username.value !== '') {
                const message = {
                    typ: "add",
                    userName: username.value,
                }
                ws.send(JSON.stringify(message));

            }
        }

        function Confirm(title, msg, $true, $false) { /*change*/
            var $content = "<div class='dialog-ovelay '>" +
                "<div class='dialog bg-gradient'><header>" +
                " <h3> " + title + " </h3> " +
                "<i class='fa fa-close'></i>" +
                "</header>" +
                "<div class='dialog-msg'>" +
                " <p> " + msg + " </p> " +
                "</div>" +
                "<footer>" +
                "<div class='controls'>" +
                " <button class='button button-danger doAction'>" + $true + "</button> " +
                " <button class='button button-default cancelAction'>" + $false + "</button> " +
                "</div>" +
                "</footer>" +
                "</div>" +
                "</div>";
            $('body').prepend($content);
            $('.doAction').click(function () {
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
                    image.value = ''
                }
                $(this).parents('.dialog-ovelay').fadeOut(500, function () {
                    $(this).remove();
                });
            });
            $('.cancelAction, .fa-close').click(function () {
                image.value = ''
                $(this).parents('.dialog-ovelay').fadeOut(500, function () {
                    $(this).remove();
                });
            });

        }
        document.getElementById('file').addEventListener('change', () => {
            Confirm('Send Image', 'Are you sure you want to send this image?', 'Yes', 'Cancel');
        })
        document.getElementById('leaveButton').addEventListener('click', () => {
            if (confirm('Are you sure you want to leave the chat?')) {
                const message = {
                    typ: "remove",
                    userName: username.value,
                    token: token
                }
                ws.send(JSON.stringify(message));
                username.value = '';
                document.getElementById('messages').innerHTML = "";
                setTimeout(location.reload(), 1000)
            }
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
        function alertMsg(messageObj) {
            if (messageObj.typ === 'error') {
                if (!window.alert('Username Already exist!')) { window.location.reload(); }
            }
        }
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
                                  <span class="mt-2 ms-1 px-3 py-2  bg-dark rounded-pill"
                                 style="font-family: Comic Sans MS, Comic Sans, cursive; font-size: 15px;color: #F3E5AB; border: 2px solid #F3E5AB;"><span
                                style="font-family:Comic Sans MS, Comic Sans, cursive ; font-size: 15px; color: #F3E5AB;"><span class="fw-bolder">${messageObj.msg}</span>
                                 </div>`
                token = messageObj.token !== "" ? messageObj.token : token
                document.getElementById("onlineUser").innerHTML = 'Online:' + messageObj.totalUser;
            }
            else if (messageObj.typ === 'txt-msg') {
                if (messageObj.userName === username.value) {
                    string = `<div class="container reciever mt-2 mb-2  bg-gradient   float-end  me-2" id="sender" style='border: 2px solid #F3E5AB;'  >
                                      <h4 class="mt-2 ms-1 fw-bolder" style="font-family: Comic Sans MS, Comic Sans, cursive; color: #F3E5AB;">${messageObj.userName}</h4>
                                      <p class="ms-1 text-break" style="font-family:Comic Sans MS, Comic Sans, cursive ; font-size: 15px; color: #F3E5AB;">${messageObj.msg}</p>
                                      </div>`;

                }
                else if ((messageObj.userName !== username.value)) {
                    string = ` <div class="container  reciever mt-2 mb-2  bg-dark bg-gradient float-start  ms-2"  style='border: 2px solid #F3E5AB;'
                        id="reciever" >
                        <h4 class="mt-2 ms-1 fw-bolder" style="font-family: Comic Sans MS, Comic Sans, cursive; color: #F3E5AB;">${messageObj.userName}</h4>
                        <p class="ms-1 text-break" style="font-family:Comic Sans MS, Comic Sans, cursive ; font-size: 15px; color: #F3E5AB;">${messageObj.msg}</p>
                        </div>`;

                }
            }
            else if (messageObj.typ === 'img-msg') {
                imgSrc = messageObj.msg
                if (messageObj.userName === username.value) {
                    string = `  <div class="container reciever mt-2 mb-2 bg-gradient float-end me-2 " id="sender" style='border: 2px solid #F3E5AB;'>
                                    <h4 class="mt-2 ms-1 fw-bolder" style="font-family: Comic Sans MS, Comic Sans, cursive; color: #F3E5AB;">${messageObj.userName}</h4>
                                    <p class=" mx-auto ms-1" ><img width="660" height="325" src=${messageObj.msg} onclick="imgClick('${messageObj.msg}')" id="imgSender" alt="" class="imgSender mx-auto"></p>
                                    </div> `;
                }
                else if ((messageObj.userName !== username.value)) {
                    string = `  <div class="container reciever mt-2 mb-2 bg-gradient float-start me-2 " id="sender" style='border: 2px solid #F3E5AB;'>
                                    <h4 class="mt-2 ms-1 fw-bolder" style="font-family: Comic Sans MS, Comic Sans, cursive; color: #F3E5AB;">${messageObj.userName}</h4>
                                    <p class=" ms-1 mx-auto" ><img  width="660" height="325" src=${messageObj.msg} onclick="imgClick('${messageObj.msg}')" id="imgSender" alt="" class="imgSender mx-auto"></p>
                                    </div> `;
                }


            }
            // converting element string to dom
            let messageElement = getElementFromString(string);
            // console.log(parameterElement);
            messages.appendChild(messageElement);

        }
    }
})

