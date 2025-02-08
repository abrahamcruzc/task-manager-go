class UI {
    constructor() {
        this.taskList = document.getElementById('taskList');
        this.loadingIndicator = document.getElementById('loadingIndicator');
        this.taskModal = new bootstrap.Modal(document.getElementById('taskModal'));
        this.deleteModal = new bootstrap.Modal(document.getElementById('deleteModal'));
        this.currentTaskId = null;
    }

    showLoading() {
        this.loadingIndicator.style.display = 'block';
        this.taskList.innerHTML = '';
    }

    hideLoading() {
        this.loadingIndicator.style.display = 'none';
    }

    getStatusBadgeClass(status) {
        const statusClasses = {
            'To do': 'bg-secondary',
            'In progress': 'bg-primary',
            'Completed': 'bg-success'
        };
        return statusClasses[status] || 'bg-secondary';
    }

    displayTasks(tasks) {
        this.hideLoading();
        this.taskList.innerHTML = '';
        
        if (tasks.length === 0) {
            this.taskList.innerHTML = `
                <div class="col-12 text-center">
                    <p class="text-muted">No tasks found. Add a new task to get started!</p>
                </div>
            `;
            return;
        }

        tasks.forEach(task => {
            const taskCard = document.createElement('div');
            taskCard.className = 'col-md-4 mb-4';
            taskCard.innerHTML = `
                <div class="card h-100">
                    <div class="card-body">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <h5 class="card-title mb-0">${task.name}</h5>
                            <span class="badge ${this.getStatusBadgeClass(task.status)}">${task.status}</span>
                        </div>
                        <p class="card-text">${task.description}</p>
                    </div>
                    <div class="card-footer bg-transparent border-top-0">
                        <div class="btn-group w-100">
                            <button class="btn btn-outline-primary" onclick="app.editTask(${task.ID})">
                                <i class="bi bi-pencil"></i> Edit
                            </button>
                            <button class="btn btn-outline-danger" onclick="app.confirmDelete(${task.ID})">
                                <i class="bi bi-trash"></i> Delete
                            </button>
                        </div>
                    </div>
                </div>
            `;
            this.taskList.appendChild(taskCard);
        });
    }

    showTaskModal(task = null) {
        const modalTitle = document.getElementById('modalTitle');
        const taskForm = document.getElementById('taskForm');
        const taskId = document.getElementById('taskId');
        const taskName = document.getElementById('taskName');
        const taskDescription = document.getElementById('taskDescription');
        const taskStatus = document.getElementById('taskStatus');

        modalTitle.textContent = task ? 'Edit Task' : 'Add New Task';
        taskId.value = task ? task.ID : '';
        taskName.value = task ? task.name : '';
        taskDescription.value = task ? task.description : '';
        taskStatus.value = task ? task.status : 'To do';

        this.taskModal.show();
    }

    showDeleteModal(taskId) {
        this.currentTaskId = taskId;
        this.deleteModal.show();
    }

    hideModals() {
        this.taskModal.hide();
        this.deleteModal.hide();
    }

    showToast(message, type = 'success') {
        const toastContainer = document.querySelector('.toast-container');
        const toast = document.createElement('div');
        toast.className = `toast align-items-center text-white bg-${type}`;
        toast.setAttribute('role', 'alert');
        toast.setAttribute('aria-live', 'assertive');
        toast.setAttribute('aria-atomic', 'true');

        toast.innerHTML = `
            <div class="d-flex">
                <div class="toast-body">${message}</div>
                <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
            </div>
        `;

        toastContainer.appendChild(toast);
        const bsToast = new bootstrap.Toast(toast);
        bsToast.show();

        toast.addEventListener('hidden.bs.toast', () => {
            toast.remove();
        });
    }

    getFormData() {
        return {
            ID: parseInt(document.getElementById('taskId').value) || undefined,
            name: document.getElementById('taskName').value,
            description: document.getElementById('taskDescription').value,
            status: document.getElementById('taskStatus').value
        };
    }
}
