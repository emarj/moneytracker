$(
    function () {
      $('#is-shared').change(function () {
        if ($('input[name="Amount"]').val() == 0) {
          $(this).prop('checked', false);
          return
        }
        var checked = $(this).prop('checked');
        $('#shared-perc').prop('disabled', !checked)
        $('#own-amount').prop('disabled', !checked)
        $('#shared-amount').prop('disabled', !checked)

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

      $('#toggle-geoloc').change(function () {
        var checked = $(this).prop('checked');
        $('#position').prop('disabled', !checked);
      });

      $('#shared-amount').change(updateFromAmount);

      $('#type-select').change(function () {
        value = $(this).find(":selected").attr('value')
        cond = (value == 1)

        if (cond) {
          $("#category-select").val('1')
        } else {
          $("#category-select").val('0')
        }
      });

      $(".clickable").click(function () {
        window.document.location = $(this).data("href");
      });

      //var wpid = navigator.geolocation.watchPosition(geo_success, geo_error, geo_options);



    });

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



  function geo_success(position) {
    if (!$('#geoloc-toggle').prop('disabled')) {
      $('#position').val(position.coords.latitude + ',' + position.coords.longitude);
    }
  }

  function geo_error() {
    console.log("Sorry, no position available.");
  }

  var geo_options = {
    enableHighAccuracy: true,
    maximumAge: 30000,
    timeout: 27000
  };