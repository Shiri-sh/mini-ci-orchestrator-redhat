const buildButton = document.getElementById("buildButt");
const repoInput = document.getElementById("repo");
const branchInput = document.getElementById("branch");
const buildsList = document.getElementById("buildsList");

buildButton.addEventListener("click", (event) => triggerBuild(event));
watchBuilds()

async function triggerBuild(event) {
  event.preventDefault();
  const repo = repoInput.value;
  const branch = branchInput.value;

  try {
    const response = await fetch("http://localhost:8080/build/create", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ repo, branch })
    })
    const newBuild = await response.json();
    console.log("Build triggered:", newBuild);
  }
  catch (error) {
    console.log("Error triggering build:", error.message);
    return;
  }
}

async function fetchBuilds() {
  try {
    console.log("Fetching builds...")
    const response = await fetch("http://localhost:8080/builds")
    const builds = await response.json()
    console.log("Fetched builds:", builds)
    setTable(builds)
  }
  catch (error) {
    console.log("Error fetching builds:", error.message)
  }
}
function setTable(builds){
  buildsList.innerHTML = builds.map(build => `
    <tr>
      <td>${build.id}</td>
      <td>${build.repo}</td>
      <td>${build.status}</td>
      <td>${new Date(build.created_at).toLocaleString()}</td>
    </tr>
    `).join("")
}
function watchBuilds() {
  console.log("Watching builds...")
  fetchBuilds()
  setInterval(fetchBuilds, 5000)
}
console.log("App initialized")