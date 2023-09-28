
const toggleThemeButton = document.getElementById("toggleThemeButton");
const themeLink = document.getElementById("theme");


function toggleTheme() {

  if (themeLink.getAttribute("href") === "style.css") {
    themeLink.setAttribute("href", "darkmode.css"); 
  } else {
    themeLink.setAttribute("href", "lightmode.css"); 
  }
}


toggleThemeButton.addEventListener("click", toggleTheme);


document.getElementById("generateButton").addEventListener("click", function () {
  const blueprintData = {
    name: document.getElementById("name").value,
    actions: [], 
    configMapNames: document.getElementById("configMapNames").value,
    containerName: document.getElementById("containerName").value,
    backupCommand: document.getElementById("backupCommand").value,
    restoreCommand: document.getElementById("restoreCommand").value
  };

  const actionsCheckboxes = document.getElementsByName("actions");
  actionsCheckboxes.forEach((checkbox) => {
    if (checkbox.checked) {
      blueprintData.actions.push(checkbox.value);
    }
  });

  console.log("Blueprint Data:", blueprintData);

  fetch("/generateBlueprint", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(blueprintData),
  })
  .then((response) => {
    console.log("Response Status:", response.status);
    return response.json();
  })
  .then((data) => {
    console.log("Response Data:", data);
    const generatedBlueprint = data.blueprintYAML
      .replace("YourNamespace", document.getElementById("namespace").value)
      .replace("YourBackupCommandHere", document.getElementById("backupCommand").value)
      .replace("YourRestoreCommandHere", document.getElementById("restoreCommand").value);

    document.getElementById("generatedBlueprint").innerText = generatedBlueprint;
  })
  .catch((error) => {
    console.error("Error:", error);
  });
});
