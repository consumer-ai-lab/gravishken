const express = require('express');
const bodyParser = require('body-parser');
const app = express();
const port = 8081;

app.use(bodyParser.json());

let messages = [];
let lastMessageId = 0;

// Route to get all messages
app.get('/messages', (req, res) => {
    res.json(messages);
});

// Long polling route to get new messages
app.get('/poll', (req, res) => {
    const lastKnownId = parseInt(req.query.lastId) || 0;
    
    const sendNewMessages = () => {
        const newMessages = messages.filter(msg => msg.id > lastKnownId);
        if (newMessages.length > 0) {
            res.json(newMessages);
        }
    };

    // Check immediately for new messages
    sendNewMessages();

    // If no new messages, wait for up to 30 seconds
    const intervalId = setInterval(sendNewMessages, 1000);
    const timeoutId = setTimeout(() => {
        clearInterval(intervalId);
        res.json([]); // Send empty array if no new messages after 30 seconds
    }, 30000);

    // Clean up interval and timeout if client disconnects
    req.on('close', () => {
        clearInterval(intervalId);
        clearTimeout(timeoutId);
    });
});

// Route to post a new message
app.post('/message', (req, res) => {
    const { text } = req.body;
    if (!text) {
        return res.status(400).json({ error: 'Message text is required' });
    }

    lastMessageId++;
    const newMessage = {
        id: lastMessageId,
        text,
        timestamp: new Date().toISOString()
    };
    messages.push(newMessage);
    res.status(201).json(newMessage);
});

// Route to clear all messages (for testing purposes)
app.delete('/messages', (req, res) => {
    messages = [];
    lastMessageId = 0;
    res.sendStatus(204);
});

app.listen(port, () => {
    console.log(`Polling server running at http://localhost:${port}/`);
});