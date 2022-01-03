videojs.registerPlugin('backForwardButtons', function() {
    let myPlayer = this,
        jumpAmount = 10,
        controlBar,
        insertBeforeNode,
        newElementBB = document.createElement('div'),
        newElementFB = document.createElement('div'),
        newImageBB = document.createElement('img'),
        newImageFB = document.createElement('img');
  
    // +++ Assign IDs for later element manipulation +++
    newElementBB.id = 'backButton';
    newElementFB.id = 'forwardButton';
  
    // +++ Assign properties to elements and assign to parents +++
    newImageBB.setAttribute('src', '/static/images/back-button.png');
    newElementBB.appendChild(newImageBB);
    newImageFB.setAttribute('src', '/static/images/forward-button.png');
    newElementFB.appendChild(newImageFB);
  
    // +++ Get control-bar and insert elements +++
    controlBar = myPlayer.$('.vjs-control-bar');
    // Get the element to insert buttons in front of in control-bar
    insertBeforeNode = myPlayer.$('.vjs-volume-panel');
  
    // Insert the button div in proper location
    controlBar.insertBefore(newElementBB, insertBeforeNode);
    controlBar.insertBefore(newElementFB, insertBeforeNode);

    // +++ Add event handlers to jump back or forward +++
    // Back button logic, don't jump to negative times
    newElementBB.addEventListener('click', function () {
      let newTime,
          rewindAmt = jumpAmount,
          videoTime = myPlayer.currentTime();
      if (videoTime >= rewindAmt) {
        newTime = videoTime - rewindAmt;
      } else {
        newTime = 0;
      }
      myPlayer.currentTime(newTime);
    });

    // Forward button logic, don't jump past the duration
    newElementFB.addEventListener('click', function () {
      let newTime,
          forwardAmt = jumpAmount,
          videoTime = myPlayer.currentTime(),
          videoDuration = myPlayer.duration();
      if (videoTime + forwardAmt <= videoDuration) {
        newTime = videoTime + forwardAmt;
      } else {
        newTime = videoDuration;
      }
      myPlayer.currentTime(newTime);
    });
  });

videojs.getPlayer('myPlayerID').ready(function () {
    let myPlayer = this;
    myPlayer.backForwardButtons();
});

let vid1 = videojs('myPlayerID');
vid1.on('dblclick', function () {
  vid1.currentTime() + 5;
});
vid1.on('dblclick', function () {
  vid1.currentTime() - 5;
});
