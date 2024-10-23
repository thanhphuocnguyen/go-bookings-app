let attention = Prompt();
(function () {
  'use strict';
  window.addEventListener(
    'load',
    function () {
      // Fetch all the forms we want to apply custom Bootstrap validation styles to
      let forms = document.getElementsByClassName('needs-validation');
      // Loop over them and prevent submission
      Array.prototype.filter.call(forms, function (form) {
        form.addEventListener(
          'submit',
          function (event) {
            if (form.checkValidity() === false) {
              event.preventDefault();
              event.stopPropagation();
            }
            form.classList.add('was-validated');
          },
          false
        );
      });
    },
    false
  );
})();

function notify(msg, msgType) {
  notie.alert({
    type: msgType,
    text: msg,
  });
}

function notifyModal(title, text, icon, confirmationButtonText) {
  Swal.fire({
    title: title,
    html: text,
    icon: icon,
    confirmButtonText: confirmationButtonText,
  });
}

function Prompt() {
  let toast = function (c) {
    const { msg = '', icon = 'success', position = 'top-end' } = c;

    const Toast = Swal.mixin({
      toast: true,
      title: msg,
      position: position,
      icon: icon,
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
      didOpen: (toast) => {
        toast.addEventListener('mouseenter', Swal.stopTimer);
        toast.addEventListener('mouseleave', Swal.resumeTimer);
      },
    });

    Toast.fire({});
  };

  let success = function (c) {
    const { msg = '', title = '', footer = '' } = c;

    Swal.fire({
      icon: 'success',
      title: title,
      text: msg,
      footer: footer,
    });
  };

  let error = function (c) {
    const { msg = '', title = '', footer = '' } = c;

    Swal.fire({
      icon: 'error',
      title: title,
      text: msg,
      footer: footer,
    });
  };

  async function custom(c) {
    const { msg = '', title = '' } = c;

    const { value: formValues } = await Swal.fire({
      title: title,
      html: msg,
      backdrop: false,
      focusConfirm: false,
      showCancelButton: true,
      willOpen: () => {
        const elem = document.getElementById('reservation-dates-modal');
        const rp = new DateRangePicker(elem, {
          format: 'yyyy-mm-dd',
          showOnFocus: true,
        });
      },
      didOpen: () => {
        document.getElementById('start').removeAttribute('disabled');
        document.getElementById('end').removeAttribute('disabled');
      },
      preConfirm: () => {
        return [
          document.getElementById('start').value,
          document.getElementById('end').value,
        ];
      },
    });

    if (formValues) {
      Swal.fire(JSON.stringify(formValues));
    }
  }

  return {
    toast: toast,
    success: success,
    error: error,
    custom: custom,
  };
}
