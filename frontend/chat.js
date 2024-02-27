var log = document.querySelector(".chat")

// function appendLog(container, msg, date) {
//     var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
//     log.prepend(container); // Prepend the container instead of appending
//     container.append(msg);
//     msg.append(date);
   
//     if (doScroll) {
//         log.scrollTop = log.scrollHeight - log.clientHeight;
//     }
// }

// 
function appendLog(container, msg, date, prepend = false) {
    var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
  
    if (prepend) {
      log.prepend(container); // Prepend the container instead of appending
    } else {
      log.appendChild(container); // Append the container
    }
    
    container.appendChild(msg);
    msg.appendChild(date);
    
    if (doScroll) {
      log.scrollTop = log.scrollHeight - log.clientHeight; // Scroll to the bottom of the chatbox
    }
  }
  


// Create messages from data and append to log container
function CreateMessages(data, currId) {
    data.reverse();

    //new addition 
    //log.innerHTML = "";
    
    // Iterate over the data in reverse order
    for (let i = data.length - 1; i >= 0; i--) {
      const { id, sender_id, content, date } = data[i];
  
      // Check if the message with the same ID already exists in the chatbox
      if (document.getElementById(`message-${id}`)) {
        continue; // Skip this message and move to the next iteration
      }
  
      // Create the message elements
      const receiverContainer = document.createElement("div");
      receiverContainer.className = sender_id == currId ? "sender-container" : "receiver-container";
      
      const receiver = document.createElement("div");
      receiver.className = sender_id == currId ? "sender" : "receiver";
      receiver.innerText = content;
  
      const messagedate = document.createElement("div");
      messagedate.className = "chat-time";
      messagedate.innerText = date.slice(0, -3);
  
      // Set a unique ID for each message element
      receiverContainer.id = `message-${id}`;
  
      // Append the message elements to the log container
      appendLog(receiverContainer, receiver, messagedate, true);
      //log.innerHTML = "";

        // Scroll to the bottom of the chat box after adding the new message
        chatBox.scrollTop = chatBox.scrollHeight;
      //firstId = id;
    }
  }
  
// Create users from data and append to users container
function sendMsg(conn, rid, msg, msg_type) {
    console.log("receiver ID from sendmsg:",rid)
    if (!conn) {
        return false;
    }
    console.log("Message input Value:", msg.value);
    if (!msg.value) {
        console.log("No message");
        return false;
    }
    console.log("Sending message to receiver ID:", rid);

    let msgData = {
        id: 0,
        sender_id: 0,
        receiver_id: rid,
        content: msg.value,
        date: '',
        msg_type: msg_type,
        is_typing: false
    }

    conn.send(JSON.stringify(msgData))
    msg.value = "";
    updateUsers()


    console.log("Message sent to receiver ID:", rid);
    //console.log("Message content:", msg.value);
    console.log("Message Data:", msgData.content);
    return false;
};

var chatBox = document.querySelector(".chat");

function OpenChat(rid, conn, data, currId, firstId) {
    document.getElementById('id' + rid).style.fontWeight = "400";

    for (var i = 0; i < unread.length; i++) {
        if (unread[i][0] == rid) {
            unread[i][1] = 0;
        }
    }

    let oldElem = document.querySelector(".send-wrapper");
    let newElem = oldElem.cloneNode(true);
    oldElem.parentNode.replaceChild(newElem, oldElem);

    document.querySelector(".chat-user-username").innerText = allUsers.filter(u => {
        return u.id == rid;
    })[0].username;
    document.querySelector(".chat-wrapper").style.display = "flex";
    var msg = document.getElementById("chat-input");

    log.innerHTML = "";
    

    var typingTimer;
    var typingInterval = 3000;
    var isTyping = false;

    document.querySelector("#chat-input").addEventListener("keydown", function(event) {
        // Add typing listener and send typing notification
        if (event.keyCode !== 13 && !isTyping) { // Ignore the Enter key and if already typing
            console.log("User started typing");
            isTyping = true;
    
            clearTimeout(typingTimer);
            typingTimer = setTimeout(function() {
                if (isTyping) {
                    console.log("User stopped typing");
                    sendTypingStatus(conn, rid, false); // Send typing status as false when stopped typing
                    isTyping = false;
                }
            }, typingInterval);
    
            sendTypingStatus(conn, rid, true); // Send typing status as true when started typing
            console.log("Sending typing status:", true);
        }
    });
    
    // Send message on button click or Enter keypress and fetch new messages
    document.querySelector("#send-btn").addEventListener("click", function() {
        sendMsg(conn, rid, msg, 'msg');
        offset = null;
        console.log("Sending message From button:", msg.value);
        firstId = firstId + 10;
        let resp = getData('http://localhost:8000/message?receiver=' + rid + '&firstId=512' + '&offset=10');
        resp.then(value => {
            console.log("enter resp then =>", value);
            console.log("Enter value", value);
            if (value && value.length > 0) {
                const lastIndex = value.length - 1;
                firstId = value[lastIndex].id;
                lastFetchedId = firstId;
                  console.log("Button lastFetchedID", firstId);
              }
            CreateMessages(value, currId);
            console.log("Sending message From button:", value);
            console.log("firstId from button ", firstId);
            // Scroll to the bottom of the chat box after adding the new message
            chatBox.scrollTop = chatBox.scrollHeight;
        }).catch();
        //resetScroll();
        //firstId = 512;
        log.innerHTML = "";
        clearTimeout(typingTimer); // Clear the previous timer
        isTyping = false; // Reset the typing flag
    });
    
    document.querySelector("#chat-input").addEventListener("keydown", function(event) {
        if (event.keyCode === 13) { // Send message on Enter keypress
            offset = 0;
            sendMsg(conn, rid, msg, 'msg');
            console.log("Sending message From Enter:", msg.value);
            firstId = firstId + 10;
            let resp = getData('http://localhost:8000/message?receiver=' + rid + '&firstId=512'  + '&offset=10');
            resp.then(value => {
                console.log("enter resp then =>", value);
                console.log("Enter value", value);
                if (value && value.length > 0) {
                    const lastIndex = value.length - 1;
                    firstId = value[lastIndex].id;
                    lastFetchedId = firstId;
                      console.log("Button lastFetchedID", firstId);
                  }
                CreateMessages(value, currId);
                console.log("Sending message From Enter:", value);
                console.log("firstId from enter", firstId);
                // Scroll to the bottom of the chat box after adding the new message
                chatBox.scrollTop = chatBox.scrollHeight;
            }).catch();
            log.innerHTML = "";
            clearTimeout(typingTimer); // Clear the previous timer
            isTyping = false; // Reset the typing flag
            offset += value.length;
        }
    });

        //var firstId = 69 ;
        var offset = 10; // Initialize the offset
        var currentScrollPos = 0; // Store the current scroll position
    
        var lastFetchedId = null; // Store the ID of the last fetched message
     
    
      
      // Create the debounced event handler
       debouncedScrollHandler = debounce(function() {
    
    
          console.log("Fetching more messages w/ scroll id",rid);
        if (chatBox.scrollTop === 0) {
          console.log("User has scrolled to the top, message ID is:", firstId);
    
          
      
          // User has scrolled to the top
          // Fetch more chat log history for the current receiver
          let resp = getData('http://localhost:8000/message?receiver=' + rid + '&firstId=' + firstId + '&offset=' + offset);
          resp.then(value => {
            // Process the retrieved chat log history
            console.log("scroll fetch messages", value);
            
            // Filter out messages with duplicate IDs
            value = value.filter(message => message.id !== lastFetchedId);
      
            if (value.length > 0) {
              const lastIndex = value.length - 1;
              firstId = value[lastIndex].id;
              lastFetchedId = firstId;
              console.log("debounce firstId", firstId);
              console.log("debounce lastFetchedId", lastFetchedId);
              console.log("debounce message log", value);
              console.log("debounce receiver ID", rid);
            }
            // Calculate current scroll position
            const currentScrollPos = chatBox.scrollHeight - chatBox.scrollTop;
            CreateMessages(value, currId);
    
            // Calculate new scroll position
          var newScrollPos = chatBox.scrollHeight - currentScrollPos;
          chatBox.scrollTop = newScrollPos;
      
            offset += value.length; // Increment the offset by 10
            //offset =+ 10;
          }).catch();
        }
      }, 300); // Adjust the debounce delay as needed (e.g., 300 milliseconds)
      
      // Attach the debounced event handler to the scroll event
         chatBox.addEventListener("scroll", debouncedScrollHandler);
    
    
            // Define the debounce function
            function debounce(func, delay) {
            let timer;
            return function() {
            clearTimeout(timer);
            timer = setTimeout(func, delay);
                };
            }
            
      

    if (data == null) {
        return;
    }

    CreateMessages(data, currId);
}


function animateDots(typingInterval) {
    var typingText = document.querySelector("#typing-text"); // Get the element with ID "typing-text"
    typingText.textContent = "is typing"; // Set the text content to "is typing"
    var dotsContainer = document.querySelector("#typing-dots"); // Get the element with ID "typing-dots"
    dotsContainer.innerHTML = "";
    var dotsCount = 1;
    var maxDots = 3;
    var dotsAnimationInterval = setInterval(function() {
        dotsContainer.innerHTML = ".".repeat(dotsCount);
        dotsCount++;
        if (dotsCount > maxDots) {
            dotsCount = 1;
        }
    }, 500);

    // Stop the dots animation after a certain duration
    setTimeout(function() {
        clearInterval(dotsAnimationInterval);
        typingText.textContent = ""; // Clear the typing notification after the animation ends
    }, typingInterval);
}

function sendTypingStatus(conn, receiverId, isTyping) {
    let typingData = {
        receiver_id: receiverId,
        is_typing: isTyping
    };

    conn.send(JSON.stringify(typingData));
    console.log("Sent typing status:", typingData);
}


// close chat
document.querySelector(".close-chat").addEventListener("click", function() {
    document.querySelector(".chat-wrapper").style.display = "none";
    console.log("closed chat");
    resetScroll();
});
// reset scroll variables and remove scroll event listener
function resetScroll() {
    var chatBox = document.querySelector(".chat");
  
    // Check if the chatBox element and debouncedScrollHandler function exist
    if (chatBox && typeof debouncedScrollHandler === "function") {
      // Remove the scroll event listener
      chatBox.removeEventListener("scroll", debouncedScrollHandler);
    }
        // Reset variables and scroll position
        firstId = 512;
        offset = 0;
        currentScrollPos = 0;
        lastFetchedId = null;
        value = null;
        closeChat = true;
}


