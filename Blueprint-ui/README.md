## Kanister Blueprint Editor

This project is a web-based editor for generating Kanister blueprints. Kanister is an open-source framework that allows you to build composable and reusable backup and restore workflows for Kubernetes applications. This editor simplifies the process of creating these blueprints through a user-friendly interface.

### Project Structure
```plaintext
│   go.mod
│   go.sum
│   go.work
│   main.go
│
├───blueprint
│       blueprint.go
│       go.mod
│       go.sum
│
├───Kanister Examples
│       BlueprintExamples.go
│       go.mod
│       go.sum
│
└───static
        darkmode.css
        editor.html
        index.html
        lightmode.css
        script.js
        style.css
```

### How It Works

1. **HTML Form (index.html):**
   - Users input Kanister blueprint details through a form.
   - The form includes fields for blueprint name, actions (backup, restore, delete), output paths, namespace, container name, image, and command details.
   - Users can toggle between dark mode and light mode using the "Toggle Dark Mode" button.

2. **JavaScript (script.js):**
   - JavaScript handles user interactions.
   - Toggles between dark and light mode styles based on user preference.
   - Sends the form data as a POST request to the `/generateBlueprint` endpoint when the "Generate Blueprint" button is clicked.
   - Displays the generated Kanister blueprint in the HTML page.

3. **Go Backend (main.go and blueprint.go):**
   - **main.go:**
     - Serves the static files (HTML, CSS, JavaScript) from the `/static/` endpoint.
     - Handles the `/generateBlueprint` endpoint where it receives the form data via a POST request.
     - Passes the received data to `generateBlueprintHandler` function in `blueprint.go`.
     - Responds with the generated Kanister blueprint in JSON format.

   - **blueprint.go:**
     - Contains the logic to generate Kanister blueprints based on the received data.
     - Defines data structures for blueprints and actions.
     - Generates YAML-formatted blueprints using the input data.
     - Contains utility functions for blueprint generation.

### How to Use
1. Run `main.go` & Open `index.html` in a web browser.
2. Fill out the form with your Kanister blueprint details.
3. Click "Generate Blueprint" to create your Kanister blueprint.
4. The generated blueprint will be displayed on the page.
5. You can toggle between dark and light mode using the "Toggle Dark Mode" button.

This project integrates HTML, JavaScript, and Go to provide a user-friendly interface for generating Kanister blueprints, making the process of creating backup and restore workflows for Kubernetes applications more accessible and efficient.

From Veeam hackathon, not complete, needed to verify blueprint examples, we also planned a second web page that uses a Go program in this repo under `/kanister Examples`, `BlueprintExamples.go`, that fetches all .yaml files in Kanister `/examples` directory that have "blueprint" in their name. These .YAML files are then presented through UI using just their names in a dropdown list, allowing the user to select an existing blueprint and edit in UI, finally saving locally to be applied.

With this, you have a custom blueprint creator page and a preset example editor page. Finally, it would have been ideal to containerize and run in K8s.
