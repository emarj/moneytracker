var isShared;
var isLocOK;

function updateFromPerc() {
  var quota = $('#shared-perc').val();
  var amount = $('input[name="Amount"]').val();
  var ownAmount = amount * (100 - quota) / 100;
  var shrAmount = amount * (quota) / 100;

  $('#own-amount').val(ownAmount);
  $('#shared-amount').val(shrAmount);
}

function updateFromAmount() {
  var shrAmount = $('#shared-amount').val();
  var amount = $('input[name="Amount"]').val();
  var quota = shrAmount * 100 / amount;

  $('#shared-perc').val(quota).prop('disabled', true);

  $('#own-amount').val(amount - shrAmount);
}

function getPosition() {
  navigator.geolocation.getCurrentPosition(geo_success, geo_error, geo_options);
}

function geo_success(position) {
    console.log(position.coords.latitude)
    $('#position').val(position.coords.latitude + ',' + position.coords.longitude);
}

function geo_error() {
  console.log("Sorry, no position available.");
}

var geo_options = {
  enableHighAccuracy: true,
  maximumAge: 30000,
  timeout: 27000
};

$(function () {

      isShared = $('#is-shared').prop('checked');
      isLocOK = $('#loc-check').prop('checked');

      if (isLocOK) {
        getPosition();
      }

      $('#loc-check').change(function () {
        isLocOK = $(this).prop('checked');
        $('#position').prop('disabled', !isLocOK);
        if (!isLocOK) {
            $('#position').val('');
        } else {
          getPosition();
        }
      });

      $('#is-shared').click(function(e){
        if ($('input[name="Amount"]').val() == 0) {
          e.preventDefault();
          return
        }
      });

      $('#is-shared').change(function () {
        
        isShared = $(this).prop('checked');

        $('#shared-perc').prop('disabled', !isShared)
        $('#own-amount').prop('disabled', !isShared)
        $('#shared-amount').prop('disabled', !isShared)

        if (checked) {
          updateFromPerc();
        } else {
          $('#shared-amount').val('');
          $('#own-amount').val('');
        }
      });

      $('input[name="Amount"]').on('input', function () {
        if ($('#is-shared').prop('checked')) {
          updateFromPerc();
        }
      });
      $('#shared-perc').on('input', updateFromPerc)
        .on('dblclick', function () {
          $(this).val(50);
          updateFromPerc();
        });

      $('#shared-amount').change(updateFromAmount);


      $(".clickable").click(function () {
        window.document.location = $(this).data("href");
      });



    });
