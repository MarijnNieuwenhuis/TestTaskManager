// Tasks Controller - Handles all task-related interactions
import { Controller } from "https://unpkg.com/@hotwired/stimulus@3.2.2/dist/stimulus.js"

export default class extends Controller {
    static targets = ["input", "error", "list", "label", "priorityInput", "taskCount"]

    // Track active filters
    activeFilters = new Set()

    // Create a new task
    async create(event) {
        event.preventDefault()

        const title = this.inputTarget.value.trim()

        if (!title) {
            this.showError("Please enter a task title")
            return
        }

        // Get selected priority and color
        const selectedInput = this.priorityInputTargets.find(input => input.checked)
        const priority = selectedInput ? selectedInput.value : "📋"
        const color = selectedInput ? selectedInput.dataset.color : "#6c757d"

        try {
            const response = await fetch("/api/tasks", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ title, priority, color }),
            })

            const data = await response.json()

            if (!response.ok) {
                this.showError(data.error || "Failed to create task")
                return
            }

            // Clear input and error
            this.inputTarget.value = ""
            this.hideError()

            // Reload page to show new task
            window.location.reload()
        } catch (error) {
            this.showError("Network error: Could not create task")
            console.error("Create task error:", error)
        }
    }

    // Toggle task completion status
    async toggle(event) {
        const taskId = this.getTaskId(event.target)
        const checkbox = event.target
        const label = checkbox.nextElementSibling

        // Optimistic UI update
        const wasChecked = checkbox.checked
        this.updateTaskUI(label, wasChecked)

        try {
            const response = await fetch(`/api/tasks/${taskId}/toggle`, {
                method: "PATCH",
            })

            if (!response.ok) {
                // Revert on error
                checkbox.checked = !wasChecked
                this.updateTaskUI(label, !wasChecked)

                const data = await response.json()
                this.showError(data.error || "Failed to toggle task")
            }
        } catch (error) {
            // Revert on network error
            checkbox.checked = !wasChecked
            this.updateTaskUI(label, !wasChecked)

            this.showError("Network error: Could not toggle task")
            console.error("Toggle task error:", error)
        }
    }

    // Delete a task
    async delete(event) {
        const taskId = this.getTaskId(event.target)

        if (!confirm("Are you sure you want to delete this task?")) {
            return
        }

        const listItem = event.target.closest("li")

        try {
            const response = await fetch(`/api/tasks/${taskId}`, {
                method: "DELETE",
            })

            if (!response.ok) {
                const data = await response.json()
                this.showError(data.error || "Failed to delete task")
                return
            }

            // Remove from UI with animation
            listItem.style.opacity = "0"
            listItem.style.transform = "translateX(-20px)"
            listItem.style.transition = "all 0.3s ease"

            setTimeout(() => {
                listItem.remove()

                // Check if list is empty
                const remainingTasks = document.querySelectorAll("li[data-task-id]")
                if (remainingTasks.length === 0) {
                    window.location.reload()
                }
            }, 300)
        } catch (error) {
            this.showError("Network error: Could not delete task")
            console.error("Delete task error:", error)
        }
    }

    // Filter tasks by priority
    filterByPriority(event) {
        const priority = event.target.dataset.priority
        const button = event.target

        // Toggle filter
        if (this.activeFilters.has(priority)) {
            this.activeFilters.delete(priority)
            button.classList.remove("active")
        } else {
            this.activeFilters.add(priority)
            button.classList.add("active")
        }

        this.applyFilters()
    }

    // Clear all filters
    clearFilters() {
        this.activeFilters.clear()

        // Remove active state from all filter buttons
        document.querySelectorAll('[data-action*="filterByPriority"]').forEach(btn => {
            btn.classList.remove("active")
        })

        this.applyFilters()
    }

    // Apply active filters to task list
    applyFilters() {
        const tasks = this.listTarget.querySelectorAll('[data-task-id]')
        let visibleCount = 0

        tasks.forEach(task => {
            const taskPriority = task.dataset.priority

            // Show if no filters active OR priority matches any active filter
            if (this.activeFilters.size === 0 || this.activeFilters.has(taskPriority)) {
                task.classList.remove('d-none')
                visibleCount++
            } else {
                task.classList.add('d-none')
            }
        })

        // Update task count
        if (this.hasTaskCountTarget) {
            const totalCount = tasks.length
            const countText = this.activeFilters.size > 0
                ? `Showing ${visibleCount} of ${totalCount} tasks`
                : `Showing ${totalCount} tasks`
            this.taskCountTarget.textContent = countText
        }
    }

    // Helper: Get task ID from element's data attribute
    getTaskId(element) {
        const listItem = element.closest("li[data-task-id]")
        return listItem.dataset.taskId
    }

    // Helper: Update task label styling
    updateTaskUI(label, isCompleted) {
        if (isCompleted) {
            label.classList.add("text-decoration-line-through", "text-muted")
        } else {
            label.classList.remove("text-decoration-line-through", "text-muted")
        }
    }

    // Helper: Show error message
    showError(message) {
        if (this.hasErrorTarget) {
            this.errorTarget.textContent = message
            this.errorTarget.classList.remove("d-none")

            // Auto-hide after 5 seconds
            setTimeout(() => this.hideError(), 5000)
        }
    }

    // Helper: Hide error message
    hideError() {
        if (this.hasErrorTarget) {
            this.errorTarget.classList.add("d-none")
            this.errorTarget.textContent = ""
        }
    }
}
