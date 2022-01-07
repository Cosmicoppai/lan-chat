let modal = document.getElementById("myModal");
let modalImg = document.getElementById("img01");
let captionText = document.getElementById("caption");
function imgClick() {
    let img = document.getElementById("imgSender");
    modal.style.display = "block";
    modalImg.src = img.src;
    captionText.innerHTML = img.alt;
}
let span = document.getElementsByClassName("close")[0];
span.onclick = function () {
    modal.style.display = "none";
}

