class TaskService {
    constructor() {
        this.baseUrl = '/tasks';
    }

    async getAllTasks() {
        try {
            const response = await fetch(this.baseUrl);
            if (!response.ok) throw new Error('Failed to fetch tasks');
            return await response.json();
        } catch (error) {
            console.error('Error fetching tasks:', error);
            throw error;
        }
    }

    async createTask(taskData) {
        try {
            const response = await fetch(this.baseUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(taskData)
            });
            if (!response.ok) {
                const error = await response.text();
                throw new Error(error || 'Failed to create task');
            }
            return await response.json();
        } catch (error) {
            console.error('Error creating task:', error);
            throw error;
        }
    }

    async updateTask(taskId, taskData) {
        try {
            const response = await fetch(`${this.baseUrl}/${taskId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(taskData)
            });
            if (!response.ok) {
                const error = await response.text();
                throw new Error(error || 'Failed to update task');
            }
            return await response.json();
        } catch (error) {
            console.error('Error updating task:', error);
            throw error;
        }
    }

    async deleteTask(taskId) {
        try {
            const response = await fetch(`${this.baseUrl}/${taskId}`, {
                method: 'DELETE'
            });
            if (!response.ok) throw new Error('Failed to delete task');
            return true;
        } catch (error) {
            console.error('Error deleting task:', error);
            throw error;
        }
    }
}
