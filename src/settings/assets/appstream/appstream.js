
function request(methood, url, data) {
    var req = new XMLHttpRequest();
    req.open(methood, url, false);
    req.send(data);
    return req.responseText
}

function Init(url, pid) {
    function onWheel(event) {
        console.log(event)
        request("POST", "https://" + url + "/runner/mouseevent/" + pid, JSON.stringify({ type: "wheel", coords: { x: event.offsetX, y: event.offsetY }, wheel_delta: { x: event.deltaX, y: window.scrollY } }))
    }

    window.onscroll = onWheel

    let canvas = document.getElementById("canvas")
    canvas.onmousedown = (event) => {
        request("POST", "https://" + url + "/runner/mouseevent/" + pid, JSON.stringify({ type: "leftDown", coords: { x: event.offsetX, y: event.offsetY } }))
    }
    canvas.onmouseup = (event) => {
        request("POST", "https://" + url + "/runner/mouseevent/" + pid, JSON.stringify({ type: "leftUp", coords: { x: event.offsetX, y: event.offsetY } }))
    }
    canvas.onclick = (event) => { }
    canvas.ondblclick = (event) => { }
    canvas.onmousemove = (event) => { }

    

    // if (canvas.addEventListener) {
    //     if ('onwheel' in document) {
    //         // IE9+, FF17+, Ch31+
    //         canvas.addEventListener("wheel", onWheel);
    //     } else if ('onmousewheel' in document) {
    //         // устаревший вариант события
    //         canvas.addEventListener("mousewheel", onWheel);
    //     } else {
    //         // Firefox < 17
    //         canvas.addEventListener("MozMousePixelScroll", onWheel);
    //     }
    // } else { // IE8-
    //     canvas.attachEvent("onmousewheel", onWheel);
    // }

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
    newImage.src = "https://" + url + "/runner/image/" + pid + "?time=" + new Date().getTime();
}

