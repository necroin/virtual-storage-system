
function request(methood, url, data) {
    var req = new XMLHttpRequest();
    req.open(methood, url, false);
    req.send(data);
    return req.responseText
}

function Init(url, pid){
    let canvas = document.getElementById("canvas")
    canvas.onclick = (event) => {
        console.log(event.offsetX);
        console.log(event.offsetY);
        request("POST", "https://"+url+"/runner/clicked/"+pid, JSON.stringify({x: event.offsetX, y:event.offsetY}))
    }

    setTimeout(window.LaunchStream, 0, url, pid)
}

function LaunchStream(url, pid) {
    let canvas = document.getElementById("canvas")
    let ctx = canvas.getContext('2d')
    newImage = new Image();
    newImage.onload = function () {
        canvas.width = newImage.width
        canvas.height = newImage.height
        ctx.drawImage(newImage, 0, 0, newImage.width, newImage.height);
        setTimeout(window.LaunchStream, 0, url, pid)
    }
    newImage.src = "https://"+url+"/runner/image/"+pid+"?time=" + new Date().getTime();
}

function MouseClicked(element) {

}