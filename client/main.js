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
    .then((polls) => displayPolls(polls))
    .catch((error) => console.error("Error loading polls: ", error));
}

function displayPolls(polls) {
  const pollList = document.getElementById("poll-list");
  pollList.innerHTML = "";

  for (const poll of polls) {
    const pollElement = document.createElement("div");
    pollElement.classList.add("poll");
    pollElement.innerHTML = `
            <h3>${poll.question} - Votes ${poll.total_votes}</h3>
            <div class="options">
                ${poll.options
                  .map(
                    (option, index) => `
                    <div class="option" data-index="${index}">
                    ${option} <span class="vote-count">${
                      poll.votes[index] || 0
                    }</span>
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
  }
}

function handleVote(buttonElement) {
  const pollID = parseInt(buttonElement.dataset.pollId);
  const optionIndex = parseInt(buttonElement.dataset.optionIndex);

  fetch("http://localhost:8080/votes", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Access-Control-Allow-Origin": "*", // Allow all origins
      Referer: window.location.href, // Add referer
    },
    body: JSON.stringify({ poll_id: pollID, option_index: optionIndex }),
  })
    .then((response) => {
      if (!response.ok) {
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
