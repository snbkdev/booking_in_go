function checkAvailability(roomID, csrfToken) {
    let html = `
        <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
            <div class="row" id="reservation-dates-modal">
                <div class="col">
                    <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
                </div>
                <div class="col">
                    <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">
                </div>
            </div>
        </form>
    `;

    attention.custom({
        msg: html,
        title: "Choose your dates",

        willOpen: () => {
            const elem = document.getElementById('reservation-dates-modal');
                const rp = new DateRangePicker(elem, {
                    format: 'yyyy-mm-dd',
                    showOnFocus: true,
                    minDate: new Date(),
            })
        },

        didOpen: () => {
            document.getElementById('start').removeAttribute('disabled');
            document.getElementById('end').removeAttribute('disabled');
        },

        callback: function(result) {
            if (!result) {
                return;
            }

            // result приходит из preConfirm: [start, end].
            // Читать значения из формы нельзя — модалка к этому моменту уже закрыта.
            let formData = new FormData();
            formData.append("start", result[0]);
            formData.append("end", result[1]);
            formData.append("csrf_token", csrfToken);
            formData.append("room_id", roomID);

            fetch('/search-availability-json', {
                method: "post",
                body: formData,
            })
            .then(response => response.json())
            .then(data => {
                if (data.ok) {
                    attention.custom({
                    icon: 'success',
                    showConfirmButton: false,
                    msg: '<p>Room is available!</p>' + '<p><a href="/book-room?id=' + data.room_id + '&s=' + data.start_date + '&e=' + data.end_date + '" class="btn btn-primary">Book now!</a></p>',
                    })
                } else {
                    attention.error({
                    msg: "No availability",})
                }
            })
        }
    });
}
