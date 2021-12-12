const firstNameElem = document.getElementById('firstName');
const lastNameElem = document.getElementById('lastName');
const emailElem = document.getElementById('email');
const passwordElem = document.getElementById('password');
const confirmPasswordElem = document.getElementById('confirmPassword');

const firstNameFieldElem = document.getElementById('firstNameField');
const lastNameFieldElem = document.getElementById('lastNameField');
const confirmPasswordFieldElem = document.getElementById(
  'confirmPasswordField'
);
const submitBtnElem = document.querySelector('button[type="submit"]');

const showErrorMessage = (errorMessage) => {
  const errorElem = document.getElementById('error-message');

  errorElem.classList.remove('is-hidden');
  errorElem.classList.add('mb-4');

  errorElem.innerHTML = `
        <p class="has-text-danger has-text-centered "> 
            ${errorMessage}
        </p>
    `;
  setTimeout(() => {
    errorElem.classList.add('is-hidden');
    errorElem.classList.remove('mb-4');
    errorElem.innerHTML = '';
  }, 5000);
};

const clearFormValues = () => {
  firstNameElem.value = '';
  lastNameElem.value = '';
  emailElem.value = '';
  passwordElem.value = '';
  confirmPasswordElem.value = '';
};

const signUp = async (firstName, lastName, email, password) => {
  try {
    const res = await fetch('/api/signup', {
      method: 'POST',
      body: JSON.stringify({ email, password, firstName, lastName }),
      headers: {
        'Content-Type': 'application/json',
      },
      withCredentials: true
    });
    const response = await res.json();
    console.log(response, 'RESPONSE!!!');

    loginTab.click();
  } catch (err) {
    console.log(err, 'error!!!');
    showErrorMessage(err.response.data.message);
  }
};

const login = async (email, password) => {
  try {
    const res = await fetch('/api/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
      headers: {
        'Content-Type': 'application/json',
      },
      withCredentials: true
    });
    const response = await res.json();
    console.log(response, 'RESPONSE!!!');

    window.location.assign("/planner")
  } catch (err) {
    showErrorMessage(err.response.data.message);
  }
};

document.getElementById('auth-form').addEventListener('submit', (e) => {
  e.preventDefault();

  let firstName;
  let lastName;
  let confirmPassword;

  if (submitBtnElem.id !== 'login') {
    firstName = firstNameElem.value;
    lastName = lastNameElem.value;
    confirmPassword = confirmPasswordElem.value;
  }

  const email = emailElem.value;
  const password = passwordElem.value;

  if (
    submitBtnElem.id !== 'login' &&
    password.trim() !== confirmPassword.trim()
  ) {
    showErrorMessage('password and confirm password do not match');
    return;
  }

  submitBtnElem.id === 'login'
    ? login(email, password)
    : signUp(firstName, lastName, email, password);
});

const signUpTab = document.getElementById('sign-up-tab');
const loginTab = document.getElementById('log-in-tab');

signUpTab.addEventListener('click', (e) => {
  signUpTab.classList.add('is-active');
  loginTab.classList.remove('is-active');

  firstNameFieldElem.classList.remove('is-hidden');
  lastNameFieldElem.classList.remove('is-hidden');
  confirmPasswordFieldElem.classList.remove('is-hidden');

  submitBtnElem.id = 'signup';

  submitBtnElem.textContent = 'Sign Up';

  clearFormValues();
});

loginTab.addEventListener('click', (e) => {
  loginTab.classList.add('is-active');
  signUpTab.classList.remove('is-active');

  firstNameFieldElem.classList.add('is-hidden');
  lastNameFieldElem.classList.add('is-hidden');
  confirmPasswordFieldElem.classList.add('is-hidden');

  submitBtnElem.id = 'login';

  submitBtnElem.textContent = 'Login';

  clearFormValues();
});
