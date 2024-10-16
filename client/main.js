document.addEventListener("DOMContentLoaded", function () {
  // Initialize and display all existing polls from the server.
  loadPolls();

  // Event listener for voting button on each poll
  document
    .getElementById("poll-list")
    .addEventListener("click", function (event) {
      if (event.target.classList.contains("vote-button")) {
        handleVote(event.target);
      }
    });
});

function loadPolls() {
  fetch("http://localhost:8080/polls")
    .then((response) => response.json())
    .then((polls) => {
      displayPolls(polls);
    })
    .catch((error) => console.error("Error loading polls: ", error));
}

function getAllSockets(polls) {
  // Logic for handling websockets (example using WebSocket API)
  for (const poll of polls) {
    const socket = new WebSocket(`ws://localhost:8080/ws/${poll.id}`); // Use the poll ID for the WebSocket connection
    socket.onopen = () => console.log("WebSocket connection established");
    socket.onmessage = (event) => {
      const data = String(event.data).split("map")[1].split("}")[0];

      function parseCustomFormat(input) {
        // Remove the brackets and split by space
        const items = input.slice(1, -1).trim().split(" ");

        // Map the items to an array of objects
        const result = items.map((item) => {
          const [key, value] = item.split(":");
          return {
            poll_id: poll.id,
            optionIndex: parseInt(key, 10), // Convert the key to an integer
            voteCount: parseInt(value, 10), // Convert the value to an integer
          };
        });

        return result;
      }

      const parsedData = parseCustomFormat(data);

      if (parsedData) {
        parsedData.forEach((data) => {
          const pollElements = document.querySelectorAll(
            `div.option[data-poll-id="${data.poll_id}"]`
          );
          console.log("111", pollElements);
          pollElements.forEach((pollElement) => {
            const voteCountElement = pollElement.querySelector(
              `.vote-count[data-poll-index="${data.optionIndex}"][data-poll-id="${data.poll_id}"]`
            );
            if (voteCountElement) {
              voteCountElement.textContent = data.voteCount;
            }
          });
        });
      }
    };
  }
}
function displayPolls(polls) {
  const pollList = document.getElementById("poll-list");
  pollList.innerHTML = ""; // Clear existing poll content

  for (const poll of polls) {
    const pollElement = document.createElement("div");
    pollElement.classList.add("poll");
    pollElement.innerHTML = `
            <h3>${poll.question}</h3>
            <div class="options">
                ${poll.options
                  .map(
                    (option, index) => `
                    <div class="option" data-poll-id="${poll.id}">
                        ${option} <span class="vote-count" data-poll-index="${index}" data-poll-id="${
                      poll.id
                    }">${0}</span>
                        <button class="vote-button" data-poll-id="${
                          poll.id
                        }" data-option-index="${index}">Vote</button>
                    </div>
                `
                  )
                  .join("")}
            </div>
        `;
    pollList.appendChild(pollElement);
    getAllSockets(polls);
  }
}

function handleVote(buttonElement) {
  const pollID = parseInt(buttonElement.dataset.pollId);
  const optionIndex = parseInt(buttonElement.dataset.optionIndex);

  fetch("http://localhost:8080/votes", {
    method: "POST",
    mode: "no-cors", // Fetch the resource with CORS disabled
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ poll_id: pollID, option_index: optionIndex }),
  })
    .then((response) => {
      if (response.status !== 200) {
        throw new Error("Voting failed");
      }
      return response.json();
    })
    .then((data) => {
      // Update vote counts on the UI for real-time display
      const voteCountElement = buttonElement
        .closest(".option")
        .querySelector(".vote-count");
      voteCountElement.textContent = parseInt(voteCountElement.textContent) + 1;
    })
    .catch((error) => console.error("Error voting: ", error));
}

function handleVoteCounts(voteCounts, pollID) {
  const pollDiv = document.querySelector(`div.poll[data-poll-id="${pollID}"]`);

  const optionElements = pollDiv.querySelectorAll(".option");
  optionElements.forEach((optionElement, index) => {
    const voteCountElement = optionElement.querySelector(".vote-count");
    voteCountElement.textContent = voteCounts[index];
  });
}
