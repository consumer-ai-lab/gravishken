(async () => {
    let port = 6200;
    let url = `http://localhost:${port}/ws`;

    let ws = new WebSocket(url);

    ws.addEventListener("open", async () => {
        console.log("open");
    });
    ws.addEventListener("message", async (m) => {
        let str = m.data;
        console.log(str);
        ws.send("\"lmao back\"")
        // let mesg = JSON.parse(m.data);
    });
})()
