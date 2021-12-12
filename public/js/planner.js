const addTaskBtn = document.getElementById('add-task');
const addTaskModal = document.getElementById('add-task-modal');

const modalCloseIcon = document.querySelector('.modal-close');

addTaskBtn.addEventListener('click', () => {
  addTaskModal.classList.add('is-active');
});

modalCloseIcon.addEventListener('click', () => {
  addTaskModal.classList.remove('is-active');
});
