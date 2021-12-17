const addTaskBtn = document.getElementById('add-task');
const addTaskModal = document.getElementById('add-task-modal');
const editTaskModal = document.getElementById('edit-task-modal');
const taskDeadlineDatePickerInput = document.querySelector(
  '#task-deadline-date-picker'
);

const datePicker = MCDatepicker.create({
  el: '#task-deadline-date-picker',
  bodyType: 'inline',
  minDate: new Date(),
});

datePicker.onSelect((date, formatedDate) => {
  taskDeadlineDatePickerInput.value = date;
  taskDeadlineDatePickerInput.checked = true;

  const taskSpecificDayText = document.getElementById('task-specific-day-text');

  taskSpecificDayText.textContent = `(${formatedDate})`;
});

datePicker.onCancel(() => {
  taskDeadlineDatePickerInput.value = '';
  taskDeadlineDatePickerInput.checked = false;
  const taskSpecificDayText = document.getElementById('task-specific-day-text');

  taskSpecificDayText.textContent = `specific day`;
});

const modalCloseIcons = document.querySelectorAll('.modal-close');

addTaskBtn.addEventListener('click', () => {
  addTaskModal.classList.add('is-active');
});

Array.from(modalCloseIcons).forEach((element) => {
  element.addEventListener('click', () => {
    addTaskModal.classList.remove('is-active');
    editTaskModal.classList.remove('is-active');
  });
});

// adding task

document
  .getElementById('add-task-form')
  .addEventListener('submit', async (e) => {
    e.preventDefault();

    const title = document.getElementById('add-task-form-title').value;
    const colorTag = document.querySelector(
      'input[name="task-color-tag"]:checked'
    ).value;
    const taskDeadline = document.querySelector(
      'input[name="task-deadline"]:checked'
    ).value;

    try {
      const res = await fetch('/api/tasks', {
        method: 'POST',
        body: JSON.stringify({
          title,
          colorTag,
          // deadline:
          //   taskDeadline.length > 8
          //     ? new Date(taskDeadline).toLocaleDateString('en-US', {
          //         year: 'numeric',
          //         month: '2-digit',
          //         day: '2-digit',
          //       })
          //     : getDate(taskDeadline),
          deadline: new Date(),
          completed: false,
        }),
      });
      const response = await res.json();
      console.log(response, 'RESPONSE!!!');
      window.location.reload();
    } catch (err) {
      console.log(err.response);
    }
  });

function getDate(value) {
  if (value === 'today') {
    return new Date(new Date().setUTCHours(23, 59, 59, 999));
  } else if (value === 'tomorrow') {
    const today = new Date();
    const tomorrow = new Date(today);
    tomorrow.setDate(tomorrow.getDate() + 1);
    return new Date(tomorrow.setUTCHours(23, 59, 59, 999)).toLocaleDateString(
      'en-US',
      { year: 'numeric', month: '2-digit', day: '2-digit' }
    );
  } else if (value === 'soon') {
    const today = new Date();
    const tomorrow = new Date(today);
    tomorrow.setDate(tomorrow.getDate() + 7);
    return new Date(tomorrow.setUTCHours(23, 59, 59, 999)).toLocaleDateString(
      'en-US',
      { year: 'numeric', month: '2-digit', day: '2-digit' }
    );
  } else if (value === 'someday') {
    const today = new Date();
    const tomorrow = new Date(today);
    tomorrow.setDate(tomorrow.getDate() + 30);
    return new Date(tomorrow.setUTCHours(23, 59, 59, 999)).toLocaleDateString(
      'en-US',
      { year: 'numeric', month: '2-digit', day: '2-digit' }
    );
  }
}

const taskContainer = document.getElementsByClassName('task-container');

Array.from(taskContainer).forEach((element) => {
  element.addEventListener('click', function (e) {
    e.stopPropagation();

    if (e.target === this) {
      element
        .querySelector('#sub-task-container')
        .classList.toggle('is-hidden');
    }
  });
});

const taskCompleteInput = document.getElementsByClassName('task-complete');

Array.from(taskContainer).forEach((element) => {
  element.addEventListener('click', async function (e) {
    e.stopPropagation();
    try {
      await fetch(`/api/tasks/${e.target.id}/completed`, {
        method: 'PATCH',
        body: JSON.stringify({
          completed: e.target.checked,
        }),
        headers: {
          'Content-Type': 'application/json',
        },
        withCredentials: true,
      });

      const titleElem = document.getElementById(`task-${e.target.id}`);

      titleElem.classList.toggle('completed');
    } catch (err) {
      console.log(err, err.response, 'this is the error');
    }
  });
});

const taskEditIcon = document.getElementsByClassName('task-edit');

Array.from(taskEditIcon).forEach((element) => {
  element.addEventListener('click', (e) => {
    e.stopPropagation();

    editTaskModal.classList.add('is-active');

    const taskId = element.getAttribute('data-task-id');
    const taskDeadline = element.getAttribute('data-task-deadline');
    const taskTitle = element.getAttribute('data-task-title');
    const taskColorTag = element.getAttribute('data-task-color-tag');

    const titleElem = editTaskModal.querySelector('#edit-task-form-title');

    console.log(taskDeadline, 'taskDeadline');

    titleElem.value = taskTitle;
    editTaskModal.setAttribute('data-task-id', taskId);
    getRadioByValue(taskColorTag).checked = true;
  });
});

Array.from(document.getElementsByClassName('task-delete')).forEach(
  (element) => {
    element.addEventListener('click', async (e) => {
      e.stopPropagation();

      const taskId = element.getAttribute('data-task-id');

      try {
        await fetch(`/api/tasks/${taskId}`, {
          method: 'DELETE',
        });

        window.location.reload();
      } catch (err) {
        console.log(err.response);
      }
    });
  }
);

document
  .getElementById('edit-task-form')
  .addEventListener('submit', async (e) => {
    e.preventDefault();

    const title = document.getElementById('edit-task-form-title').value;
    const colorTag = document.querySelector(
      'input[name="task-color-tag-edit"]:checked'
    ).value;

    const taskId = editTaskModal.getAttribute('data-task-id');
    // const taskDeadline = document.querySelector(
    //   'input[name="task-deadline"]:checked'
    // ).value;

    try {
      await fetch(`/api/tasks/${taskId}/edits`, {
        method: 'PATCH',
        body: JSON.stringify({
          title,
          colorTag,
          // deadline:
          //   taskDeadline.length > 8
          //     ? new Date(taskDeadline).toLocaleDateString('en-US', {
          //         year: 'numeric',
          //         month: '2-digit',
          //         day: '2-digit',
          //       })
          //     : getDate(taskDeadline),
          deadline: new Date(),
        }),
      });

      window.location.reload();
    } catch (err) {
      console.log(err.response);
    }
  });

function getRadioByValue(v) {
  var inputs = editTaskModal.querySelectorAll('#task-color-tag');
  console.log(
    Array.from(inputs).map((input) => input.value),
    'inputs!!!',
    v
  );
  for (var i = 0; i < inputs.length; i++) {
    console.log(inputs[i].value, v, 'inputs[i].value == v');
    if (inputs[i].type == 'radio' && inputs[i].value == v) {
      return inputs[i];
    }
  }
  return false;
}
