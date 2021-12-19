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

Array.from(taskCompleteInput).forEach((element) => {
  element.addEventListener('change', async function (e) {
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

Array.from(document.getElementsByClassName('sub-task-delete')).forEach(
  (element) => {
    element.addEventListener('click', async (e) => {
      e.stopPropagation();

      const taskId = element.getAttribute('data-task-id');
      const subTaskId = element.getAttribute('data-sub-task-id');

      await fetch(`/api/subtask/${subTaskId}/${taskId}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
        withCredentials: true,
      });

      const deletedElement = document.getElementById(`sub-task-${subTaskId}`);

      console.log(deletedElement, 'this is the deleted element!!!');

      if (deletedElement.parentNode) {
        deletedElement.parentNode.removeChild(deletedElement);
      }
    });
  }
);

Array.from(document.getElementsByClassName('sub-input')).forEach((element) => {
  element.addEventListener('keypress', async (e) => {
    if (e.keyCode == 13) {
      const taskId = e.target.getAttribute('data-task-id');

      if (!e.target.value) return;

      // send post request...
      const res = await fetch(`/api/subtask/${taskId}`, {
        method: 'POST',
        body: JSON.stringify({
          completed: false,
          title: e.target.value,
        }),
        headers: {
          'Content-Type': 'application/json',
        },
        withCredentials: true,
      });

      const response = await res.json();

      // TODO: get the right id in the response and attach to created subtask element

      console.log(response.data, 'response.data!!!');

      const subTasksElem = document.getElementById(`sub-tasks-${taskId}`);
      const subTaskElem = document.createElement('div');

      subTaskElem.id = 'sub-task-';

      subTaskElem.innerHTML = `
      

                    <div
                class="
                  is-flex is-justify-content-space-between is-align-items-center 
                  my-2
                "
                style="width: 50%;"
              >
                <div class="pretty p-svg p-curve">
                  <input type="checkbox" />
                  <div class="state p-primary">
                    <!-- svg path -->
                    <svg class="svg svg-icon" viewBox="0 0 20 20">
                      <path
                        d="M7.629,14.566c0.125,0.125,0.291,0.188,0.456,0.188c0.164,0,0.329-0.062,0.456-0.188l8.219-8.221c0.252-0.252,0.252-0.659,0-0.911c-0.252-0.252-0.659-0.252-0.911,0l-7.764,7.763L4.152,9.267c-0.252-0.251-0.66-0.251-0.911,0c-0.252,0.252-0.252,0.66,0,0.911L7.629,14.566z"
                        style="stroke: white; fill: white"
                      ></path>
                    </svg>
                    <label> ${e.target.value} </label>
                  </div>
                </div>
                <span class="material-icons is-size-6"> delete </span>
              </div>
      
      `;

      subTasksElem.append(subTaskElem);

      e.target.value = '';
    } else {
      return false;
    }
  });
});

Array.from(document.getElementsByClassName('sub-task-complete')).forEach(
  (element) => {
    element.addEventListener('click', async (e) => {
      const subTaskId = element.getAttribute('data-sub-task-id');
      const taskId = element.getAttribute('data-task-id');

      await fetch(`/api/subtask/${subTaskId}/${taskId}`, {
        method: 'PATCH',
        body: JSON.stringify({
          completed: e.target.checked,
        }),
        headers: {
          'Content-Type': 'application/json',
        },
        withCredentials: true,
      });

      const titleElem = document.getElementById(`sub-task-title-${subTaskId}`);

      console.log(titleElem, 'title element');

      titleElem.classList.toggle('completed');

      console.log('---after---');

      console.log(document.getElementById(`sub-task-${subTaskId}`));
    });
  }
);

Array.from(document.getElementsByClassName('task-delete')).forEach(
  (element) => {
    element.addEventListener('click', async (e) => {
      e.stopPropagation();

      const taskId = element.getAttribute('data-task-id');

      Swal.fire({
        title: 'Are you sure?',
        text: "You won't be able to revert this!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonText: 'Yes, delete it!',
        cancelButtonText: 'No, cancel!',
        reverseButtons: true,
      }).then(async (result) => {
        if (result.isConfirmed) {
          try {
            await fetch(`/api/tasks/${taskId}`, {
              method: 'DELETE',
            });
            Swal.fire('Deleted!', 'Your file has been deleted.', 'success');
            setTimeout(() => {
              window.location.reload();
            }, 1000);
          } catch (err) {
            console.log(err.response);
          }
        } else if (
          /* Read more about handling dismissals below */
          result.dismiss === Swal.DismissReason.cancel
        ) {
          Swal.fire('Cancelled', 'You cancelled', 'error');
        }
      });
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
