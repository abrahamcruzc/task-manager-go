class App {
    constructor() {
        this.taskService = new TaskService();
        this.ui = new UI();
        this.initialize();
    }

    async initialize() {
        // Event Listeners
        document.getElementById('addTaskBtn').addEventListener('click', () => this.ui.showTaskModal());
        document.getElementById('saveTaskBtn').addEventListener('click', () => this.saveTask());
        document.getElementById('confirmDeleteBtn').addEventListener('click', () => this.deleteTask());

        // Load initial tasks
        await this.loadTasks();
    }

    async loadTasks() {
        try {
            this.ui.showLoading();
            const tasks = await this.taskService.getAllTasks();
            this.ui.displayTasks(tasks);
        } catch (error) {
            this.ui.showToast('Failed to load tasks', 'danger');
            this.ui.hideLoading();
        }
    }

    async saveTask() {
        try {
            const taskData = this.ui.getFormData();
            if (taskData.ID) {
                await this.taskService.updateTask(taskData.ID, taskData);
                this.ui.showToast('Task updated successfully');
            } else {
                await this.taskService.createTask(taskData);
                this.ui.showToast('Task created successfully');
            }
            this.ui.hideModals();
            await this.loadTasks();
        } catch (error) {
            this.ui.showToast(error.message, 'danger');
        }
    }

    async editTask(taskId) {
        try {
            const tasks = await this.taskService.getAllTasks();
            const task = tasks.find(t => t.ID === taskId);
            if (task) {
                this.ui.showTaskModal(task);
            }
        } catch (error) {
            this.ui.showToast('Failed to load task details', 'danger');
        }
    }

    confirmDelete(taskId) {
        this.ui.showDeleteModal(taskId);
    }

    async deleteTask() {
        try {
            await this.taskService.deleteTask(this.ui.currentTaskId);
            this.ui.showToast('Task deleted successfully');
            this.ui.hideModals();
            await this.loadTasks();
        } catch (error) {
            this.ui.showToast('Failed to delete task', 'danger');
        }
    }
}

// Initialize the application
const app = new App();
