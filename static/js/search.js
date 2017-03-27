    function addDirection(elem) {
        console.log(elem);
        var directionField = $.parseHTML(`
            <form action="#" onSubmit="event.preventDefault(); return searchDirection(this);">
                <input type="text" />
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
        queryElem = $(elem).children()[0]
        console.log($(queryElem).val());
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
