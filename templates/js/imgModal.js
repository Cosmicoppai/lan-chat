let modal = document.getElementById("myModal");
let img = document.getElementById("imgSender");
let img2 = document.getElementById("imgSender2");
let modalImg = document.getElementById("img01");
let captionText = document.getElementById("caption");
// img.onclick = function(){
//   modal.style.display = "block";
//   modalImg.src = this.src;
//   captionText.innerHTML = this.alt;
// }
// img2.onclick = function(){
//   modal.style.display = "block";
//   modalImg.src = this.src;
//   captionText.innerHTML = this.alt;
// }
// span.onclick = function() {
//   modal.style.display = "none";
// }
img.addEventListener('click', function () {
  modal.style.display = "block";
  modalImg.src = this.src;
  captionText.innerHTML = this.alt;
});
let span = document.getElementsByClassName("close")[0];
span.addEventListener('click', function () {
  modal.style.display = "none";
});
img2.addEventListener('click', function () {
  modal.style.display = "block";
  modalImg.src = this.src;
  captionText.innerHTML = this.alt;
});
