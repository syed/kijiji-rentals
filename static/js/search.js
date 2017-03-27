    function addDirection(elem) {
        console.log(elem);
        var directionField = $.parseHTML(`
            <form action="#" onSubmit="event.preventDefault(); return searchDirection(this);">
                <input type="text" />
                <select id="mode" onchange="calcRoute();">
                  <option value="DRIVING">Driving</option>
                  <option value="WALKING">Walking</option>
                  <option value="BICYCLING">Bicycling</option>
                  <option value="TRANSIT">Transit</option>
                </select>
                <a href="#" onclick="removeDirection(this)">
                <i class="fa fa-trash" style="color:red"></i>
              </a>
            </form>`);

        $(elem).parent().append(directionField);
    }

    function addPlace(elem) {
        console.log(elem);
        var placeField = $.parseHTML(`
            <form action="#" onSubmit="event.preventDefault(); return searchPlace(this);">
                <input type="text" />
                <a href="#" onclick="removePlace(this)">
                <i class="fa fa-trash" style="color:red"></i>
              </a>
            </form>`);

        $(elem).parent().append(placeField);
    }

    function searchDirection(elem) {
        var queryElem = $(elem).children()[0]
        var destStr = $(queryElem).val();

        var travelModeSelectEleme= $(elem).children()[1]
        var travelMode = $(travelModeSelectEleme).val();


        var lat = parseFloat($(elem).parent().parent().find("input[name='lat']").val());
        var lng = parseFloat($(elem).parent().parent().find("input[name='lng']").val());
        var originLatLng = { lat:lat, lng: lng};

        addDirectionOverlay(originLatLng, destStr, travelMode);

        return false;
    }

    function searchPlace(elem) {
        queryElem = $(elem).children()[0]
        console.log($(queryElem).val());
        return false;
    }

    function removeDirection(elem) {
        $(elem).parent().remove();
    }

    function removePlace(elem) {
        $(elem).parent().remove();
    }
