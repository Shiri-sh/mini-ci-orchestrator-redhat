const buildButton = document.getElementById("buildButt")
const repoInput = document.getElementById("repo")
const branchInput = document.getElementById("branch")

buildButton.addEventListener("click", triggerBuild)
async function triggerBuild() {

  const repo = repoInput.value
  const branch = branchInput.value

  await fetch("http://localhost:8080/build", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ repo, branch })
  })

  alert("Build started")
}