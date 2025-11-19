// Stimulus.js Application Bootstrap
import { Application } from "https://unpkg.com/@hotwired/stimulus@3.2.2/dist/stimulus.js"
import TasksController from "./controllers/tasks_controller.js"

// Initialize Stimulus application
window.Stimulus = Application.start()

// Configure Stimulus development experience
Stimulus.debug = true
Stimulus.warnings = true

// Register controllers
Stimulus.register("tasks", TasksController)

console.log("Stimulus application loaded")
