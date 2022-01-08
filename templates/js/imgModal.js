let modal = document.getElementById("myModal");
let modalImg = document.getElementById("img01");
let captionText = document.getElementById("caption");
function imgClick(imgSrc) {
    modal.style.display = "block";
    modalImg.src = imgSrc;
}
let span = document.getElementsByClassName("close")[0];
span.onclick = function () {
    modal.style.display = "none";
}

