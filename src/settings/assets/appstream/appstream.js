
function Init(url, pid){
    setTimeout(window.LaunchStream, 0, url, pid)
}

function LaunchStream(url, pid) {
    let canvas = document.getElementById("canvas")
    let ctx = canvas.getContext('2d')
    newImage = new Image();
    newImage.onload = function () {
        canvas.width = newImage.width
        canvas.height = newImage.height
        ctx.drawImage(newImage, 0, 0, newImage.width, newImage.height, 0,0,1600,900);
        setTimeout(window.LaunchStream, 0, url, pid)
    }
    newImage.src = "https://"+url+"/runner/image/"+pid+"?time=" + new Date().getTime();
}