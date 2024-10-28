document
  .getElementById('check-availability-button')
  .addEventListener('click', function () {
    let html = `
        <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
            <div class="form-row">
                <div class="col">
                    <div class="form-row" id="reservation-dates-modal">
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
                        </div>
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">
                        </div>

                    </div>
                </div>
            </div>
        </form>
    `;
    attention.custom({
      title: 'Choose your dates',
      msg: html,
      willOpen: () => {
        const elem = document.getElementById('reservation-dates-modal');
        new DateRangePicker(elem, {
          format: 'yyyy-mm-dd',
          showOnFocus: true,
        });
      },
      didOpen: () => {
        document.getElementById('start').removeAttribute('disabled');
        document.getElementById('end').removeAttribute('disabled');
      },
      callback: () => {
        let form = document.getElementById('check-availability-form');
        let formData = new FormData(form);
        formData.append('csrf_token', '{{.CSRFToken}}');
        fetch('/search-availability-json', {
          method: 'post',
          body: formData,
        })
          .then((response) => {
            return response.json();
          })
          .then((data) => {
            console.log(data);
          });
      },
    });
  });
